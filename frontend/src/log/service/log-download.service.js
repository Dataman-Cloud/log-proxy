/**
 * Created by lixiaoyan on 2017/4/12.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('logDownload', logDownload);

    /* @ngInject */

    function logDownload(){
        return{
            downloadFile: downloadFile
        };

        function downloadFile(filename, data) {
            var logMessage = [];
            var success = false;
            var contentType = 'text/plain;charset=utf-8';
            angular.forEach(data,function(value, kay){
                this.push(value.message)
            },logMessage);

            try {
                // Try using msSaveBlob if supported
                var blob = new Blob([logMessage], {type: contentType});
                if (navigator.msSaveBlob) {
                    navigator.msSaveBlob(blob, filename);
                } else {
                    // Try using other saveBlob implementations, if available
                    var saveBlob = navigator.webkitSaveBlob || navigator.mozSaveBlob || navigator.saveBlob;
                    if (saveBlob === undefined) throw "Not supported";
                    saveBlob(blob, filename);
                }
                success = true;
            } catch (ex) {
                console.log("saveBlob method failed with the following exception:");
                console.log(ex);
            }

            if(!success){
                // Get the blob url creator
                var urlCreator = window.URL || window.webkitURL || window.mozURL || window.msURL;
                if(urlCreator){
                    // Try to use a download link
                    var link = document.createElement('a');
                    if('download' in link){
                        // Try to simulate a click
                        try{
                            // Prepare a blob URL
                            var blob = new Blob([logMessage], { type: contentType });
                            var url = urlCreator.createObjectURL(blob);
                            link.setAttribute('href', url);

                            // Set the download attribute (Supported in Chrome 14+ / Firefox 20+)
                            link.setAttribute("download", filename);

                            // Simulate clicking the download link
                            var event = document.createEvent('MouseEvents');
                            event.initMouseEvent('click', true, true, window, 1, 0, 0, 0, 0, false, false, false, false, 0, null);
                            link.dispatchEvent(event);
                            console.log("Download link method with simulated click succeeded");
                            success = true;
                        } catch(ex) {
                            console.log("Download link method with simulated click failed with the following exception:");
                            console.log(ex);
                        }
                    }
                    if(!success){
                        // Fallback to window.location method
                        try{
                            // Prepare a blob URL
                            // Use application/octet-stream when using window.location to force download
                            console.log("Trying download link method with window.location ...");
                            var blob = new Blob([logMessage], { type: octetStreamMime });
                            var url = urlCreator.createObjectURL(blob);
                            window.location = url;
                            console.log("Download link method with window.location succeeded");
                            success = true;
                        }  catch(ex) {
                            console.log("Download link method with window.location failed with the following exception:");
                            console.log(ex);
                        }
                    }
                }
            }
            if(!success){
                // Fallback to window.open method
                console.log("No methods worked for saving the arraybuffer, using last resort window.open.  Not Implemented");
                //window.open(httpPath, '_blank', '');
            }
        }
    }
})();