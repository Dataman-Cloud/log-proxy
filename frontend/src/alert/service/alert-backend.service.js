/**
 * Created by lixiaoyan on 2016/12/2.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('alertBackend', alertBackend);

    /* @ngInject */
    function alertBackend($resource) {
        return {
            alerts: alerts,
            alert: alert
        };

        function alerts() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/alert', null, {
                'update': {method: 'PUT'}
            });
        }

        function alert(id) {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/alert/:id', {id: id});
        }
    }
})();