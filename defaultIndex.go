package main

const defaultIndex string = `<!DOCTYPE HTML>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="initial-scale=1,width=device-width">
	<title>{{ .Path }}</title>
	<style>
	  body {
			display: flex;
			flex-direction: column;
			min-height: calc(100vh - 2rem);
			margin: 0;
			padding: 1rem;
			font-size: 16px;
		}

		h1 {
			margin-top: 0;
			padding: 0 .75rem;
		}

		table {
			border-collapse: collapse;
			font-size: .875rem;
		}

		th {
			background-color: #f0f0f0;
		}

		th, td {
			padding: .5rem .75rem;
			border: 1px #dddddd solid;
		}

		footer {
			display: flex;
			justify-content: center;
			align-items: flex-end;
			flex-grow: 1;
			padding: 1rem .75rem 0 .75rem;
		}
	</style>
</head>
<body>
	<h1>{{ .Path }}</h1>
	<table>
	  <thead>
		  <tr>
			  <th>Filename</th>
			  <th>Size</th>
			  <th>Modification Time</th>
			</tr>
		</thead>
		<tbody>
		{{range .Entries}}
		  <tr>
			  <td><a href="{{ .Path }}">{{ .Name }}</a></td>
				<td>{{ .Size }}</td>
				<td>{{ .ModTime }}</td>
			</tr>
		{{end}}
		</tbody>
	</table>
	<footer>
	  This
	</footer>
</body>
</html>`
