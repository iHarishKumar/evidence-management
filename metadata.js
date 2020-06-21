const exiftool = require("exiftool-vendored").exiftool
const crypto = require('crypto-js')
const random = require('./app/mongodb/mongodb.js')
const fs = require('fs')
const mongo = require('mongodb')


var extractMetadata = function(args, imageURL, funcCallback, key, rest){
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
    var metadataEncrypt = crypto.AES.encrypt(JSON.stringify(metadata), key).toString()
    var fileEncrypt = crypto.AES.encrypt(JSON.stringify(file), key).toString()
    var metadataHash = crypto.SHA256(JSON.stringify(metadata)).toString()
    var fileHash = crypto.SHA256(JSON.stringify(file)).toString()
    // console.log(metadataHash)
    // console.log(fileHash)
    // console.log("--------------------------------")
    // console.log(arguments)
    // console.log("--------------------------------")
    
    args.push(fileEncrypt)
    args.push(metadataEncrypt)
    args.push(fileHash)
    args.push(metadataHash)

    
    var readStream = fs.createReadStream(imageURL)
    readStream.pipe(random.bucket.openUploadStream(imageURL))
    .on('finish', function(db) {
      console.log("----------------DBConnection-------------------")
      console.log(db)
      var fileArray = []
      // Now hash the file content obtained from the query to fs.chunks and fs.files
      random.db.collection(random.chunksCollection).find({files_id: db['_id']})
        .toArray(async function(err, val) {
          if(err){
            return console.log('Something went wrong while querying MongoDB.', err)
          }
          val.forEach(function(v){
            fileArray.push(v)
          })
          var hash = crypto.SHA256(JSON.stringify(fileArray)).toString()
          var dbCollectionEncrypt = crypto.AES.encrypt(JSON.stringify(db), key).toString()
          args.push(hash)
          args[4] = dbCollectionEncrypt // Need to change this. For now replacing the image url with the db entry returned after uploading the file
          let message = await funcCallback(rest[0], rest[1], rest[2], rest[3], args, rest[4], rest[5]);
          rest[6].send(message)
        })
    })
    })    
  .catch(err => console.error("Something terrible happened: ", err))
}

var decryptMetadata = function(encryptData, key){
  return crypto.AES.decrypt(encryptData, key).toString(crypto.enc.Utf8)
}

var verify = async function(res, val, key){
  var fileURL = val['DBCOLLECTION' ]
  //var readStream = fs.createReadStream(fileURL)
  var fileDataHash
  var dbCollection = decryptMetadata(fileURL, key)
  let chunkValue
  try{
    chunkValue = JSON.parse(dbCollection)
  }
  catch(err){
    console.log('Data was tampered on blockchain')
  }
  var fileArray = []
  t = new mongo.ObjectId(chunkValue['_id'])

  var ex = random.db.collection(random.chunksCollection).find({files_id: t})
    .toArray(function(err, value) {
      if(err){
        return console.log('Something went wrong while querying MongoDB', err)
      }
      value.forEach(function(v){
        fileArray.push(v)
      })
      var hashData = val[(val['Document_Data_HASH'] == undefined ? 'FIR_DATA_HASH' : 'Document_Data_HASH')]
      var dbHash = crypto.SHA256(JSON.stringify(fileArray)).toString()
      if(hashData === dbHash){
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