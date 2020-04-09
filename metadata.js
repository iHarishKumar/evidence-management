const exiftool = require("exiftool-vendored").exiftool
const crypto = require('crypto-js')
const random = require('crypto')
const key = random.randomBytes(32)

var extractMetadata = function(args, imageURL, token, funcCallback, rest){
	console.log(typeof args)
	var imageURL = args[args.length-1]
  exiftool
  .read(imageURL)
  .then(async (tags) =>{
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
    var metadataHash = crypto.AES.encrypt(JSON.stringify(metadata), "abc").toString()
    var fileHash = crypto.AES.encrypt(JSON.stringify(file), "abc").toString()
    console.log(metadataHash)
    console.log(fileHash)
    console.log("--------------------------------")
    console.log(arguments)
    console.log("--------------------------------")
    args.push(metadataHash)
    args.push(fileHash)
    let message = await funcCallback(rest[0], rest[1], rest[2], rest[3], args, rest[4], rest[5]);
    rest[6].send(message)
  })
  .catch(err => console.error("Something terrible happened: ", err))
  console.log("-----------")
}

var decryptMetadata = function(fileHash, metadataHash){
  return [crypto.AES.decrypt(fileHash, "abc").toString(crypto.enc.Utf8), crypto.AES.decrypt(metadataHash, "abc").toString(crypto.enc.Utf8)]
}

exports.extractMetadata = extractMetadata
exports.decryptMetadata = decryptMetadata