/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorAlertCtrl', MonitorAlertCtrl);
    /* @ngInject */
    
    function MonitorAlertCtrl(rules){
        var self = this;

        self.rules = rules.data.rules;
        activate();

        function activate() {

        }

    }
})();