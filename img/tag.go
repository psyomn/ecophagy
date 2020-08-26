/*
Package img is just a wrapper to use around whatever can extract and
set exit comment tags.
*/
package img

import (
	"bytes"
	"fmt"
	"log"

	"github.com/dsoprea/go-exif"
	jpgstruct "github.com/dsoprea/go-jpeg-image-structure"
)

func GetExifComment(image []byte) (string, error) {
	return "", nil
}

func SetExifComment(image []byte, comment string) []byte {
	var ret []byte
	return ret
}

func setTags(data []byte) (error, []byte) {
	// TODO: this should prpbably be factored out, as messing with
	//   tags is going to become very common
	parser := jpgstruct.NewJpegMediaParser()
	mediaContextInterface, err := parser.ParseBytes(data)
	if err != nil {
		log.Println("could not parse jpeg", err)
		return err, nil
	}

	segList := mediaContextInterface.(*jpgstruct.SegmentList)

	fmt.Println("segment list : ", segList)

	rootIb, err := segList.ConstructExifBuilder()
	if err != nil {
		log.Println("could not create exif builder: ", err)
		return err, nil
	}

	index, err := rootIb.Find(exif.IfdExifId)
	if err != nil {
		log.Println("could not find ifdexifid", err)
		return err, nil
	}

	exifBt := rootIb.Tags()[index]
	exifIb := exifBt.Value().Ib()

	uc := exif.TagUnknownType_9298_UserComment{
		EncodingType:  exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
		EncodingBytes: []byte("TEST COMMENT"),
	}

	err = exifIb.SetStandardWithName("UserComment", uc)
	if err != nil {
		log.Println("could not set standard with name 'usercomment'", err)
		return err, nil
	}

	_, s, err := segList.FindExif()
	if err != nil {
		log.Println("could not findexif", err)
		return err, nil
	}

	err = s.SetExif(rootIb)
	if err != nil {
		log.Println("could not set exif in rootib", err)
		return err, nil
	}

	b := new(bytes.Buffer)
	err = segList.Write(b)
	if err != nil {
		log.Println("could not write new buffer", err)
		return err, nil
	}

	return nil, []byte{}
}
