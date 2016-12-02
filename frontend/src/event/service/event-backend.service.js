/**
 * Created by lixiaoyan on 2016/12/2.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('eventBackend', eventBackend);

    /* @ngInject */
    function eventBackend($resource) {
        return {
            monitor: monitor
        };

        function monitor() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/prometheus');
        }
    }
})();