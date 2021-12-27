# yt-dlp-archive-viewer

This is a small, simple web frontend for YouTube archives created by [yt-dlp](https://github.com/yt-dlp/yt-dlp). 

It supports displaying thumbnails and metadata, including comments, like/dislike counts and descriptions of videos.

An example yt-dlp command that dumps all of this info is this: ``yt-dlp --write-comments --write-info-json --write-thumbnail --extractor-args "youtube:max_comments=200;comment_sort=top;max_comment_depth=1" --download-archive .downloaded "https://www.youtube.com/watch?v=JfeRBm7t0QY"``

Obviously replace the last argument with whatever you want to download.

# Usage

Either build and use the executable directly or use the Docker iamge. The latter is *greatly* recommended.

yt-dlp-archive-viewer runs on port 8000 by default.

Sample Docker command: ``docker run -d --restart unless-stopped -v /your/archive/directory:/archive privateger.docker.scarf.sh/privateger/yt-dlp-archive-viewer:latest``
