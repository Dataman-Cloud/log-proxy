(function () {
    'use strict';
    angular.module('app')
        .factory('appMonitorBackend', appMonitorBackend);
    /* @ngInject */
    function appMonitorBackend($resource) {
        //////
        return {
            listApp: listApp
        };

        function listApp(appid) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/applications', {appid: appid});
        }
    }
})();
