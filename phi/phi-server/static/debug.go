package static

// A Semi proof of concept if we were to do a very, very, light http
// service. This is all pretty gnarly, and I don't recommend you write
// this way unless if you're doing quick prototypes as I am.
const DebugPage = `
<html>
  <head>
    <script type="text/javascript">
      function statusFn() {
          var xhr = new XMLHttpRequest();
          xhr.onreadystatechange = function() {
              if (this.readyState != 4) return;
              if (this.status == 200) {
                  var data = JSON.parse(this.responseText);
                  document.getElementById("status").innerHTML = data["status"];
                  document.getElementById("version").innerHTML = data["version"];
              }
              else {
                  document.getElementById("status").innerHTML = "error";
              }
          }
          xhr.open("GET", "http://127.0.0.1:9876/status", true);
          xhr.send()
      }

      function loginFn() {
          var xhr = new XMLHttpRequest();

          xhr.onreadystatechange = function() {
              if (this.readyState != 4) return;
              if (this.status == 200) {
                  var data = JSON.parse(this.responseText);
                  document.getElementById("token").innerHTML = data["token"];
              }
              else {
                  document.getElementById("token").innerHTML = "error";
              }
          };
          xhr.open("POST", "http://127.0.0.1:9876/login", true);
          xhr.setRequestHeader('Content-Type', 'application/json');
          xhr.send(JSON.stringify({
              username: "someusername",
              password: "somesupersecretpassword"
          }));
      }

      function uploadFn() {
          var reader = new FileReader();
          var file = document.getElementById('upload_file').files[0];
          let bytes = [];

          reader.readAsArrayBuffer(file);
          reader.onload = function(e){
              var arrayBuffer = e.target.result;
              bytes = new Uint8Array(arrayBuffer);
              console.log("from onload: " + bytes.length);

              var token = document.getElementById("token").innerHTML;

              var xhr = new XMLHttpRequest();
              xhr.open('POST', "http://127.0.0.1:9876/upload/newfilename.jpg/" + (+ new Date()).toString());

              xhr.setRequestHeader("Authorization", "token " + token);
              xhr.setRequestHeader("Content-Type", "application/octet-stream");

              xhr.send(bytes);
          }
      }
    </script>
  </head>

  <body>
    <h1> Test Rest Endpoints </h1>

    <div> <h2> Check Status </h2>
      <input type="button" onclick="statusFn()" value ="check status"/>
      <table>
        <tr>
          <td> status </td> <td id="status"></td>
        </tr>
        <tr>
          <td> version </td>  <td id="version"></td>
        </tr>
      </table>
    </div>

    <div> <h2> Login </h2>
    <p> <input type="button" onclick="loginFn()" value ="login"/> <div id="login-status"></div> </p>
    <p> Token: <span id="token"></span> </p>
    </div>

    <div> <h3> Upload </h3>
    <p> <input id="upload_file" type="file" /> </p>
    <p> <input type="button" value="upload" onclick="uploadFn()" /> </p>
    </div>
  </body>

</html>
`
