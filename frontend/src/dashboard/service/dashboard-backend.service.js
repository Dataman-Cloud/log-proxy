(function () {
    'use strict';
    angular.module('app')
        .factory('dashboardBackend', dashboardBackend);
    /* @ngInject */
    function dashboardBackend($resource) {
        return {
            info: info,
            listNode: listNode
        };

        function info(data) {
            var paramObj = data || {};
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/info', {
                cluster: paramObj.cluster,
                user: paramObj.user,
                app: paramObj.app
            });
        }

        function listNode() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/nodes');
        }
    }
})();
