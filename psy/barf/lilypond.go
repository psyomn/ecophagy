/*
Package barf contains code relevant to barfing.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package barf

import (
	"errors"
	"os"

	"github.com/psyomn/ecophagy/psy/common"
)

var (
	ErrNeedSongName = errors.New("provide a song name")
)

func lilypond(args common.RunParams) common.RunReturn {
	if len(args) == 0 {
		return ErrNeedSongName
	}

	songName := args[0]

	const songContents = `\version "2.18.2"
#(set-global-staff-size 16)

\header {
  title = "song title here"
  subtitle = "subtitle here"
  composer = "Simon Symeonidis"
}

lower = \relative c {
  \clef bass
  \time 4/4

  <a e>1
  r1
  <a e>1
  r1
}

upper = \relative c'' {
  \clef treble
  \time 4/4

  r4 <a e c>4 <a e b>2
  r4 <a e c>4 <c, e g>2
}


\score {
<<
  \new ChordNames \with {
    midiInstrument = "pad 2 (warm)"
    midiMinimumVolume = #0.0
    midiMaximumVolume = #0.0
  }
  {
    \chordmode {
      c2:6 a2:min
      c2:6 a2:min
      c2:6 a2:min
      c2:6 a2:min

      c1
      c1

      d1:min
      d1:min
    }
  }
  \new Staff \with {
    instrumentName = #"5str Bass"
    midiInstrument = #"electric bass (finger)"
  }
  {
    \clef bass
    \time 4/4
  }

  \new Staff \with {
    instrumentName = #"Elec Gtr (jazz)"
    midiInstrument = #"electric guitar (jazz)"
  }
  {
    \clef treble
    \time 4/4
  }

  \new PianoStaff \with {
    instrumentName = "Hammond"
    midiInstrument = "rock organ"
  }
  <<
    \new Staff = "upper" \upper
    \new Staff = "lower" \lower
  >>

  \new DrumStaff \with {
    instrumentName = #"Drums"
  }
  {
    \drummode {
        \repeat unfold 3 { <hh bd>16 hh hhho hh }
        <hh bd>16 hhho hhho hhho

        <cymch bd hh>8 hh <sn bd hh> hh <bd hh> hh <hh sn bd> hh |
        <bd hh>8 hh <sn bd hh> hh <bd hh> hh <hh sn bd> hh |
        <cymch bd hh>8 hh <sn bd hh> hh <bd hh> hh <hh sn bd> hh |
        <bd hh>8 hh <sn bd hh> hh <bd hh> hh <hh sn bd> hh |

    }
  }
  >>

  \layout { }

  \midi {
    \tempo 4 = 100
  }
}`

	file, err := os.Create(songName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(songContents); err != nil {
		return err
	}

	return nil
}
