var  mongodb = require('mongodb')
var fs = require('fs')
var gridstore = require('mongodb').GridS

const url = 'mongodb://localhost:27017'

const dbName = 'test'
chunksCollection = 'fs.chunks'
filesCollection = 'fs.files'

let db
let bucket

mongodb.MongoClient.connect(url, { useNewUrlParser: true}, (err, client) => {
    if(err){
        return console.log(err)
    }

    db = client.db(dbName)
    bucket = new mongodb.GridFSBucket(db)
    console.log("Connect to MongoDB: ", url)
    console.log("DBName: ", dbName)
    //console.log("Bucket: ", bucket)
    //console.log("DB: ", db)
    // db.collection("fs.chunks").find({files_id: new mongodb.ObjectID('5ebead301cf0e609b9f78bc4')}).toArray(function(err, docs){
    //     console.log(docs.toLocaleString())
    //     console.log(JSON.stringfy(docs))
    //     docs.forEach(function(e){
    //         console.log(e)
    //         console.log(JSON.stringify)
    //     })
    // })
    // console.log("-----------",ex)
    // ex.each(function(err, doc){
    //     console.log("====================",doc)
    // })
    //db.fs.files.findOne('5ebbebdae124f60651c85395')
    


    // fs.createReadStream('/Users/harishgunjalli/Downloads/Docker.dmg')
    // .pipe(bucket.openUploadStream('/Users/harishgunjalli/Downloads/Docker.dmg'))
    //    .on('error', (error) => {
    //        return console.log(error)
    //    })
    //    .on('finish', (val, ex) => {
    //        console.log('Done!')
    //        console.log(val)
    //        console.log(ex)
    //        process.exit(0)
    //    })
    exports.db = db
    exports.bucket = bucket

})

exports.dbName = dbName
exports.chunksCollection = chunksCollection
exports.filesCollection = filesCollection
