(function () {
    'use strict';
    angular.module('app')
        .factory('logContextModal', logContextModal);

    /* @ngInject */
    function logContextModal($mdDialog) {

        return {
            open: open
        };

        function open(templateUrl, ev, data) {
            var dialog = $mdDialog.show({
                controller: LogContextCtrl,
                controllerAs: 'vm',
                templateUrl: templateUrl,
                parent: angular.element(document.body),
                targetEvent: ev,
                clickOutsideToClose: true,
                locals: {log: data}

            });
            return dialog;
        }

        /* @ngInject */
        function LogContextCtrl($mdDialog, log, logBackend) {
            var self = this;
            self.logContexts = [];

            activate();

            function activate() {
                logBackend.logContext({appid: log.appid, taskid: log.taskid, path: log.path, offset: log.offset})
                    .get(function (data) {
                        self.logContexts = data.data;
                    })
            }

            self.ok = function () {
                $mdDialog.hide();
            };

            self.cancel = function () {
                $mdDialog.cancel();
            };
        }
    }
})();