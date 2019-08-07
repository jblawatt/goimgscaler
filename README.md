# goimgscaler
Microservice to scale images and create thumbnails

## QueryParams

QueryParam|Name|Desc
---|---|---
`f`|Filename|Image Name
`h`|Height|Dest Image Height
`w`|Width|Dest Image Width
`m`|Method| 0 = Resize ; 1 = Fill ; 2 = Fit
`a`|Anchor| 0 = Center ; ...
`i`|Interpolation| ...

### Example
```
http://localhost:8080/?f=demo.jpg&h=100&m=0
```

## Options 
### Method
Option|Name
---|---
0|Resize
1|Fill
2|Fit

### Anchor

Option|Name
---|---
0|Center
1|TopLeft
2|Top
3|TopRight
4|Left
5|Right
6|BottomLeft
7|Bottom
8|BottomRight

### Interpolation Filter

Option|Name
---|---
0|NearestNeighbor
1|Box
2|Linear
3|Hermite
4|MitchellNetravali
5|CatmullRom
6|BSpline
7|Gaussian
8|Bartlett
9|Lanczos
10|Hann
11|Hamming
12|Blackman
13|Welch
14|Cosine


# Config File

Optional Config File: `config.yaml`
```yaml

bind: 127.0.0.1:8080

cache_dir: _cache
image_dir: input

default_filter: 0
default_method: 0
default_anchor: 0

```


# TODO: 
* [ ] Cors
* [ ] Request Domain validation 
* [ ] Safe image path
* [ ] Validate Options 