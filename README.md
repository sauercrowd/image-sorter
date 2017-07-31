# Image Sorter

Sorts stuff into a bunch of folders

Please be careful, backup your data before using it!
This is beta and I don't recommend it for using it productive environments.

I'm not responsible for any data loss.

## Example
You got a folder with jpg's at `~/Pictures` and want to sort it into `~/Sorted~`

Execute

```
image-sorter -exif -in ~/Pictures -out ~/Sorted
```

to

- find every file in ~/Pictures
- try to read and parse exif data
- create a folder with the year as name in ~/Sorted
- move the file there

If exif data couldn't be parsed or found, this file won't be moved at all

## Usage
There are two modes, either exif mode or regex
When in exif mode, it will try to use the exif data of the file for the new directory name.

If you got e.g. `.png`s which should be sorted there's not exif data, but probably your files look like that:

```
PICTURE-2016-02-01 11-01-02.png
PICTURE-2017-02-25 00-01-02.png
PICTURE-2017-05-05 00-01-02.png
```

and you want it to be sorted by its year, execute 

```
image-sorter -in ~/Pictures -out ~/Sorted -regex 'PICTURE-XXXX(\-[0-9][0-9])* [0-9]*(\-[0-9]*)*\.png' -placeholder XXXX

2017/07/31 17:54:15 File: ~/Pictures/PICTURE-2016-02-01 11-01-02.png in dir Sorted/2016/PICTURE-2016-02-01 11-01-02.png
2017/07/31 17:54:15 File: ~/Pictures/PICTURE-2017-02-25 00-01-02.png in dir Sorted/2017/PICTURE-2017-02-25 00-01-02.png
2017/07/31 17:54:15 File: ~/Pictures/PICTURE-2017-05-05 00-01-02.png in dir Sorted/2017/PICTURE-2017-05-05 00-01-02.png

```

which will sort your images based on the placeholder into

```
.
├── 2016
│   └── PICTURE-2016-02-01\ 11-01-02.png
└── 2017
    ├── PICTURE-2017-02-25\ 00-01-02.png
    └── PICTURE-2017-05-05\ 00-01-02.png
```

When using `-exif` mode, the regex will only be used for finding files, so the placeholder option won't be used at all. (default regex is `.*`)
If you wan't to find your files by regex, make sure you set regex and placeholder accordingly.