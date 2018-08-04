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

## Configuration and Setup

The token can be used to secure the connection between the minio server and the
gpd. Only events that are done via the token URL are allowed and processed.

### 1. Enable webhooks in minio

The URL needs to be the URL of the gpd server. Example:

    "webhook": {
        "1": {
                "enable": true,
                "endpoint": "https://<host>:3443/gpd/events?token=<SOME-RANDOM-TOKEN>"
        }

### 2. Create bucket and enable events on the bucket via minio command

    $ minio-client mb <server>/<bucket> --region=<region>
    $ minio-client events add <server>/<bucket> arn:minio:sqs:<region>:1:webhook --suffix .pdf

For example:

    $ minio-client mb s3/test --region=<us-east-1>
    $ minio-client events add s3/test arn:minio:sqs:us-east-1:1:webhook --suffix .pdf

The notifications can also be restricted to a certain prefix using `--prefix`.

### 3. Start the gpd service

    $ ./gpd -address :3443 -token <SOME-RANDOM-TOKEN> -cert cert.pem -key key.pem"
