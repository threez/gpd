# gpd

Ghostscript PDF Daemon is a small minio companion application.
[Minio](https://minio.io) supports bucket notifications e.g. a webhook.
This application is similar to the [thumbnailer](https://github.com/minio/thumbnailer) but applies only to PDF.

Purpose of the Ghostscript PDF Daemon is to react on webhook bucket
notifications, download the file if it is a PDF, extract the text
content into a JSON file and create thumbnails for all pages of the
document.

## Architecture

                         +-----------+
                  +------>   minio   |
                  |      +---+-------+
                  |          |
        2. down-  |          | 1. webhook
        load      |          |
                  |      +---v-------+    +-------------+
        4. upload +------+    gpd    +----> ghostscript |
        json +           +-----------+    +-------------+
        thumbnails
                                    3. analyse

## Created documents layout

The created documents have the following format.

Given the input PDF with 3 pages:

    path/in/minio/XYZ.pdf

The following files will be produced:

    path/in/minio/XYZ.json
    path/in/minio/XYZ-1.jpeg
    path/in/minio/XYZ-2.jpeg
    path/in/minio/XYZ-3.jpeg

## License

The 2-Clause BSD License.

