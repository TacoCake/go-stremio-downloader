### go-stremio-downloader
A simple go program to download Stremio series from a given package.

## Usage

To use the downloader, you need to provide the path to the package and the output directory where the files will be saved.
> **Note**: The output directory is a _prefix_, individual series will be saved in subdirectories with of their respective names.

```bash
go run . -p /path/to/package -o /output/directory
```

### Example

Let's assume the `/home/user/OnePaceStremio` directory contains a stremio package containing the series `One Pace`.

Running the following command will result in the episodes being downloaded into `/output/directory/One Pace/Season X/...`:
```bash
go run . -p /home/user/OnePaceStremio -o /output/directory
```

## License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
