package main

var templateIndex = `
<!doctype html>
<html>
    <head>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/3.0.3/normalize.min.css" type="text/css" />
        <style>/*! normalize.css v3.0.3 | MIT License | github.com/necolas/normalize.css */html{font-family:sans-serif;-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%}body{margin:0}article,aside,details,figcaption,figure,footer,header,hgroup,main,menuav,section,summary{display:block}audio,canvas,progress,video{display:inline-block;vertical-align:baseline}audio:not([controls]){display:none;height:0}[hidden],template{display:none}a{background-color:transparent}a:active,a:hover{outline:0}abbr[title]{border-bottom:1px dotted}b,strong{font-weight:bold}dfn{font-style:italic}h1{font-size:2em;margin:.67em 0}mark{background:#ff0;color:#000}small{font-size:80%}sub,sup{font-size:75%;line-height:0;position:relative;vertical-align:baseline}sup{top:-0.5em}sub{bottom:-0.25em}img{border:0}svg:not(:root){overflow:hidden}figure{margin:1em 40px}hr{box-sizing:content-box;height:0}pre{overflow:auto}code,kbd,pre,samp{font-family:monospace,monospace;font-size:1em}button,input,optgroup,select,textarea{color:inherit;font:inherit;margin:0}button{overflow:visible}button,select{text-transform:none}button,html input[type="button"],input[type="reset"],input[type="submit"]{-webkit-appearance:button;cursor:pointer}button[disabled],html input[disabled]{cursor:default}button::-moz-focus-inner,input::-moz-focus-inner{border:0;padding:0}input{line-height:normal}input[type="checkbox"],input[type="radio"]{box-sizing:border-box;padding:0}input[type="number"]::-webkit-inner-spin-button,input[type="number"]::-webkit-outer-spin-button{height:auto}input[type="search"]{-webkit-appearance:textfield;box-sizing:content-box}input[type="search"]::-webkit-search-cancel-button,input[type="search"]::-webkit-search-decoration{-webkit-appearance:none}fieldset{border:1px solid silver;margin:0 2px;padding:.35em .625em .75em}legend{border:0;padding:0}textarea{overflow:auto}optgroup{font-weight:bold}table{border-collapse:collapse;border-spacing:0}td,th{padding:0}</style>
        <style>
            html {
                font-family: "Lucida Console", Monaco, monospace;
                font-size: 12px;
                height: 100%;
            }
            body, table {
				min-height: 100%;
			}
            table {
				margin-bottom: 30px;
			}
            thead {
				border-bottom: 1px solid black;
			}
            tr {
				width: 100%;
				padding-top: 10px;
			}
			tr:nth-of-type(odd) {
				background-color: #f9f9f9;
			}
			td {
				padding: 10px 5px 0px 5px;
			}
			td:first-child {
				text-align: right;
				vertical-align: top;
				background-color: #CCC;
			}
			td:nth-child(2) {
				width: 100%;
				padding-left: 10px;
				padding-right: 0px;
			}
        </style>
    </head>
<body>
	<table>
    	<thead>
        	<tr>
                <td>file</td>
                <td>line</td>
			</tr>
		</thead>
        <tbody id="output">
		</tbody>
	</table>
    <script>
        var socket = new WebSocket("ws://{{.Host}}/ws");

        socket.onmessage = function(event) {
            obj = JSON.parse(event.data);
            output = document.getElementById('output');
            row = document.createElement("tr");
            cellFile = document.createElement("td");
            cellFile.appendChild(document.createTextNode(obj.file));
            row.appendChild(cellFile);

            cellLine = document.createElement("td");
            cellLine.appendChild(document.createTextNode(obj.text));
            row.appendChild(cellLine);

            output.appendChild(row);
            y = document.height;
            window.scroll(0,y);
        }
    </script>
</body>
</html>
`
