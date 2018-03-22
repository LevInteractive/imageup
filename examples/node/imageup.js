/**
 * An interface for Image Up microservice.
 *
 * @see https://github.com/LevInteractive/imageup/
 *
 * @module utils/imageup
 */

const request = require("request");
const config = require("config");
const fs = require("fs");

const IU_SERVER = `${config.imageup.host}:${config.imageup.port}`;

/**
 * Send a photo to imageup. You'll most likely want to use a data stream so the
 * image is never stored on disk (or if it is, it's only temporary).
 *
 * @async
 * @param {string|object} imageSrc A path to the image on disk OR a stream.
 * @param {array} sizes
 *   @param {string} sizes.name A url friendly name.
 *   @param {int} sizes.width
 *   @param {int} sizes.height
 *   @param {boolean} sizes.fit If true, will crop to size.
 * @return {array}
 */
exports.upload = function upload(imageSrc, sizes) {
  return new Promise((resolve, reject) => {
    console.info(`Uploading image via image up: ${imageSrc}`);
    request.post(
      IU_SERVER,
      {
        formData: {
          sizes: JSON.stringify(sizes),
          file:
            typeof imageSrc === "string"
              ? fs.createReadStream(imageSrc)
              : imageSrc
        }
      },
      (err, res, body) => {
        if (err) {
          reject(err);
        } else {
          resolve(JSON.parse(body));
        }
      }
    );
  });
};

/**
 * Remove any number of files from the cloud.
 *
 * @async
 * @param {array} fileNames  An array of files to remove.
 * @return {object}
 */
exports.remove = function remove(fileNames) {
  return new Promise((resolve, reject) => {
    console.info(`Removing image(s) ${fileNames}`);
    request.del(
      IU_SERVER,
      {
        formData: {
          files: fileNames.join(",")
        }
      },
      (err, res, body) => {
        if (err) {
          reject(err);
        } else {
          resolve(JSON.parse(body));
        }
      }
    );
  });
};
