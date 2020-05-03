#!/usr/bin/python

HELP = """
Usage:
  python openapi-json-to-html.py -i <inputfile> -t <title> -o <outputfile>

Example:
  python openapi-json-to-html.py -i openapi.json -t "Nuveo Code Challenge" -o index.html
"""

import sys, getopt, json

def main(argv):
  title = ''
  inputfile = ''
  outputfile = ''

  try:
    opts, args = getopt.getopt(argv,"h:i:t:o:")

  except getopt.GetoptError:
    print(HELP)
    sys.exit(2)

  for opt, arg in opts:
    if opt == '-h':
      print(HELP)
      sys.exit()
    elif opt in ("-i"):
      inputfile = arg
    elif opt in ("-t"):
      title = arg
    elif opt in ("-o"):
      outputfile = arg

  spec = json.load(open(inputfile, 'r', encoding="utf-8"))

  TEMPLATE = """
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="UTF-8">
      <title>%s</title>
      <link rel="stylesheet" type="text/css" href="./libraries/swagger-ui/3.23.8/swagger-ui.css" >
      <style>
        html
        {
          box-sizing: border-box;
          overflow: -moz-scrollbars-vertical;
          overflow-y: scroll;
        }

        *,
        *:before,
        *:after
        {
          box-sizing: inherit;
        }

        body
        {
          margin:0;
          background: #fafafa;
        }
      </style>
    </head>

    <body>
      <div id="swagger-ui"></div>
      <script src="./libraries/swagger-ui/3.23.8/swagger-ui-bundle.js"> </script>
      <script src="./libraries/swagger-ui/3.23.8/swagger-ui-standalone-preset.js"> </script>
      <script src="./libraries/swagger-ui/3.23.8/swagger-ui.js"> </script>
      <script>
      window.onload = function() {

        var spec = %s;

        const HideTopbarPlugin = function() {
          return {
            components: {
              Topbar: function() {
                return null
              }
            }
          }
        }

        const DisableTryItOutPlugin = function() {
          return {
            statePlugins: {
              spec: {
                wrapSelectors: {
                  allowTryItOutFor: () => () => false
                }
              }
            }
          }
        }

        const ui = SwaggerUIBundle({
          spec: spec,
          dom_id: '#swagger-ui',
          <!-- withCredentials: true, -->
          deepLinking: true,
          presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
          ],
          plugins: [
            SwaggerUIBundle.plugins.DownloadUrl,
            HideTopbarPlugin
          ],
          layout: "StandaloneLayout"
        })

        window.ui = ui
      }
    </script>
    </body>
  </html>
  """

  page = open(outputfile,"w+", encoding="utf-8")

  page.write(TEMPLATE % (title, json.dumps(spec)))

if __name__ == "__main__":
  main(sys.argv[1:])
