/**
 * Created by my9074 on 16/3/2.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('logcurd', logcurd);

    /* @ngInject */
    function logcurd(logContextModal) {
        return {
            openLogContext: openLogContext
        };

        function openLogContext(ev, log) {
            logContextModal.open('/src/log/modals/log-context.html', ev, log)
        }
    }
})();
