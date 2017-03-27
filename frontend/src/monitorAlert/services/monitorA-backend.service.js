/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('monitorAlertBackend', monitorAlertBackend);

    /* @ngInject */
    
    function monitorAlertBackend($resource){

        return {
            clusters : clusters,
            apps : apps,
            indicators : indicators,
            rules :rules
        }


        function clusters() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters');
        }

        function apps() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/monitor/clusters/work/apps');
        }

        function indicators() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/indicators');
        }

        function rules() {
            return $resource(BACKEND_URL_BASE.defaultBase + '/v1/alert/rules');
        }
        
    }
})();