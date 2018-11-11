const supertest  = require('supertest');
const chai = require('chai');
const expect = chai.expect;
const AdmZip = require('adm-zip');
const binaryParser = require('superagent-binary-parser');

const host = 'localhost:' + process.env.PORT;
const request = supertest(`http://${host}`);

describe('Image processing', function () {
    it('should return 400 when no file in json present', function (done) {
        request.post('/accept')
            .send({
                transformations: [{
                    "type": "crop",
                    "parameters": {
                        "x": 0,
                        "y": 0,
                        "width": 10,
                        "height": 200
                    }
                }, {
                    "type": "rotate",
                    "parameters": {
                        "angle": 1
                    }
                }]
            })
            .expect(400)
            .end(done);
    });

    it('should return 400 when there was an error downloading image', function (done) {
        request.post('/accept')
            .send({
                file: 'https://mor=mkaciuba.com/media/gallery/42661/07/Magda-6159-2_fe4618fa-ca63-11e8-ae67-0242ac120003_gallery_big1300.jpeg',
                transformations: [{
                    "type": "crop",
                    "parameters": {
                        "x": 0,
                        "y": 0,
                        "width": 10,
                        "height": 200
                    }
                }, {
                    "type": "rotate",
                    "parameters": {
                        "angle": 1
                    }
                }]
            })
            .expect(400)
            .end(done);
    });

    it('should return 400 when not image given', function (done) {
        request.post('/accept')
            .send({
                file: 'https://mort.mkaciuba.com/assets/js/mkaciuba2018v5.min.js',
                transformations: [{
                    "type": "crop",
                    "parameters": {
                        "x": 0,
                        "y": 0,
                        "width": 10,
                        "height": 200
                    }
                }, {
                    "type": "rotate",
                    "parameters": {
                        "angle": 1
                    }
                }]
            })
            .expect(400)
            .end(done);
    });
    it('should create zip file with 2 images', function (done) {
        this.timeout(5000);
        request.post('/accept')
            .send({
                file: 'https://mort.mkaciuba.com/media/gallery/42661/07/Magda-6159-2_fe4618fa-ca63-11e8-ae67-0242ac120003_gallery_big1300.jpeg',
                transformations: [{
                    "type": "crop",
                    "parameters": {
                        "x": 0,
                        "y": 0,
                        "width": 10,
                        "height": 200
                    }
                }, {
                    "type": "rotate",
                    "parameters": {
                        "angle": 1
                }
            }]
            })
            .expect(200)
            .expect('Content-Type', 'application/zip')
            .parse(binaryParser)

            .end(function(err, res) {
                if (err) {
                    return done(err);
                }

                expect(res.body.length).to.eql(151540);
                const zip = new AdmZip( res.body );
                const zipEntries = zip.getEntries();

                expect(zipEntries.length).to.eql(2);

                expect(zipEntries[0].name).to.eql('cropx-0y-0width-10height-200.jpeg');
                expect(zipEntries[1].name).to.eql('rotate-1.jpeg');

                done();
            });
    });

});