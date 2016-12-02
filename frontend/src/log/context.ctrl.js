(function () {
    'use strict';
    angular.module('app')
        .controller('LogContextCtrl', LogContextCtrl);
    /* @ngInject */
    function LogContextCtrl(logBackend, $stateParams, $anchorScroll, $location) {
        var self = this;
        var newHash = 0;
        var offset = 2;

        self.contextQueryObj = {
            appid: $stateParams.appid,
            taskid: $stateParams.taskid,
            path: $stateParams.path,
            offset: parseInt($stateParams.offset)
        };
        self.logContexts = [];

        activate();

        function activate() {
            getLogContexts();
        }

        function getLogContexts() {
            logBackend.logContext(self.contextQueryObj)
                .get(function (data) {
                    self.logContexts = data.data;
                    setTimeout(initAnchorScroll, 0);
                })
        }

        function initAnchorScroll() {
            //count location hash
            angular.forEach(self.logContexts, function (item, index) {
                if (item.offset == $stateParams.offset) {
                    if (index - offset >= 0) {
                        newHash = index - offset
                    } else {
                        newHash = index
                    }
                }
            });

            $location.hash(newHash);
            $anchorScroll();
        }
    }
})();