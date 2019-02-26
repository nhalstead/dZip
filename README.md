# dZip
- [x] FAST
- [x] SIMPLE
- [x] CROSS PLATFORM
- [x] Dependency-less

---

## Zip across all of the Operating Systems!
#### `Windows` Yes
#### `Mac` Yes, ZIP ALL THE THINGZ!
#### `Linux` ZIPP ZIPP ZIPP!!
#### `Unix` ZIP-A-DEE-DOO-DAH!!

---

|Flag|Description|
|:-:|:-|
|-file|Zip File to Target|
|-zip|If not provided it will extract `-zip` to the Zip FIle's Name|
|`...`| All extra params are treated as filenames to make a new Zip when `-zip` is passed in!|

---

#### Why does this exist?
>
> Well glad you asked, I have notived that between all of the platforms,
>  Linux, Mac and Windows the commands act different for zipping a file,
>  or not at all in windows case (without extra tools). This tool allows
>  you to add a zip function to all platforms that way its all the same,
>  simple and easy to use.
>
> This also brings me to the next thing, When using ZipArchive in PHP it requires lots
>  of addons, and takes more time than its worth but on the other hand, `exec` does not.
> Simple enough to write a script to just run the command to package a list of files at
>  execution. Dont ask why I dont just write the applicaiton in golang, its got lots
>  of movoing parts!
>
