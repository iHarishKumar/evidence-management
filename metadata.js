const exiftool = require("exiftool-vendored").exiftool
const crypto = require('crypto-js')
const random = require('crypto')
const fs = require('fs')


var extractMetadata = function(args, imageURL, funcCallback, key, rest){
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
      'fileWidth': tags.ImageWidth,
      'Latitude': tags.GPSLatitude,
      'Longtitude': tags.GPSLongitude
    }
    var file = {
      'sourceFile': tags.SourceFile,
      'directory': tags.Directory,
      'filename': tags.FileName
    }
    tags.RawDataByteOrder
    var metadataEncrypt = crypto.AES.encrypt(JSON.stringify(metadata), "abc").toString()
    var fileEncrypt = crypto.AES.encrypt(JSON.stringify(file), "abc").toString()
    var metadataHash = crypto.SHA256(JSON.stringify(metadata)).toString()
    var fileHash = crypto.SHA256(JSON.stringify(file)).toString()
    console.log(metadataHash)
    console.log(fileHash)
    console.log("--------------------------------")
    console.log(arguments)
    console.log("--------------------------------")
    
    args.push(fileEncrypt)
    args.push(metadataEncrypt)
    args.push(fileHash)
    args.push(metadataHash)

    var fileDataHash;
    var readStream = fs.createReadStream(imageURL)
    readStream.on('data', async function(chunk) {
      fileDataHash += crypto.SHA256(chunk)
      //args.push(fileDataHash)
    }).on('end', async function() {
      args.push(fileDataHash)
      let message = await funcCallback(rest[0], rest[1], rest[2], rest[3], args, rest[4], rest[5]);
      rest[6].send(message)
    })
    
    //args.push(fileDataHash)
  })
  .catch(err => console.error("Something terrible happened: ", err))
  console.log("-----------")
}

var decryptMetadata = function(encryptData, key){
  return crypto.AES.decrypt(encryptData, key).toString(crypto.enc.Utf8)
}

var verify = async function(res, val){
  var fileURL = val[(val['FIR_Image_URL'] == undefined ? 'Document_URL' : 'FIR_Image_URL')]
  var readStream = fs.createReadStream(fileURL)
  var fileDataHash
  readStream.on('data', async function(chunk) {
    fileDataHash += crypto.SHA256(chunk)
  }).on('end', async function() {
    var hashData = val[(val['Document_Data_HASH'] == undefined ? 'FIR_DATA_HASH' : 'Document_Data_HASH')]
    if(hashData == fileDataHash){
      val['isTampered'] = false
    }
    else{
      val['isTampered'] = true
    }
    res.send(val)
  })
}

exports.extractMetadata = extractMetadata
exports.decryptMetadata = decryptMetadata
exports.verify = verify