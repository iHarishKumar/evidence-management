const exiftool = require("exiftool-vendored").exiftool
const crypto = require('crypto-js')

var extractMetadata = async function(args, imageURL, token){
	console.log(typeof args)
	var imageURL = args[args.length-1]
  exiftool
    .read(imageURL)
    .then((tags) =>{
      var metadata = {
        'profileDateTime': tags.ProfileDateTime.toString(),
        'fileModified': tags.FileModifyDate.toString(),
        'filesize': tags.FileSize,
        'fileHeight': tags.ImageHeight,
        'fileWidth': tags.ImageWidth
      }
      var file = {
        'sourceFile': tags.SourceFile,
        'directory': tags.Directory,
        'filename': tags.FileName
      }
      exiftool.end()
      var metadataHash = crypto.AES.encrypt(JSON.stringify(metadata), token.split(" ")[1]).toString()
      var fileHash = crypto.AES.encrypt(JSON.stringify(file), token.split(" ")[1]).toString()
      console.log(metadataHash)
      console.log(fileHash)
      args.push(metadataHash)
      args.push(fileHash)
      return args
    })
    .catch(err => console.error("Something terrible happened: ", err))

    console.log("-----------")
}

exports.extractMetadata = extractMetadata