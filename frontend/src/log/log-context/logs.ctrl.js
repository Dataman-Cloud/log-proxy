(function () {
    'use strict';
    angular.module('app')
        .controller('LogContextCtrl', LogContextCtrl);
    /* @ngInject */
    function LogContextCtrl(logBackend, $stateParams, $timeout, logDownload) {
        var self = this;

        self.contextQueryObj = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app,
            task: $stateParams.task,
            source: $stateParams.source,
            offset: parseInt($stateParams.offset),
            page: 1,
            size: 100
        };
        self.logContexts = [];
        self.offsetIndex = 0;

        //Infinite Scroll setting
        self.infiniteLogs = {
            numLoaded_: 0,
            toLoad_: 0,
            items: [],

            // Required.
            getItemAtIndex: function (index) {
                if (index > this.numLoaded_) {
                    this.fetchMoreItems_(index);
                    return null;
                }
                return this.items[index];
            },

            // Required.
            getLength: function () {
                return this.numLoaded_ + 5;
            },

            fetchMoreItems_: function (index) {
                if (this.toLoad_ < index) {
                    this.toLoad_ += 100;
                    logBackend.logContext(self.contextQueryObj).get(angular.bind(this, function (data) {
                        if (!data.data) {
                            this.numLoaded_ = this.toLoad_ - 100;
                        } else {
                            self.contextQueryObj.page += 1;
                            this.items = this.items.concat(data.data);
                            $timeout(checkContextIndex, 0);
                            this.numLoaded_ = this.toLoad_;
                        }
                    }));
                }
            }
        };

        self.downloadFile = downloadFile;

        function downloadFile(logs) {
            self.fileName = logs[0].clusterid+'-'+logs[0].appid;
            logDownload.downloadFile(self.fileName, logs);
        }
        activate();

        function activate() {
        }

        function checkContextIndex() {
            angular.forEach(self.infiniteLogs.items, function (log, index) {
                if (log.offset == $stateParams.offset) {
                    self.offsetIndex = index;
                }
            });

        }
    }
})();
