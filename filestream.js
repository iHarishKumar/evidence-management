var crypto = require('crypto');
var fs = require('fs');

var readStream = fs.createReadStream('/Users/harishgunjalli/Desktop/2018-19/Screenshot 2019-08-27 at 1.06.57 PM.png');
var hash = crypto.createHash('sha1');
readStream
  .on('data', function (chunk) {
    hash.update(chunk);
  })
  .on('end', function () {
    console.log(hash.digest('hex'));
  });

  console.log("ergqer========")

