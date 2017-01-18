(function () {
    'use strict';
    angular.module('app')
        .factory('alertModal', alertModal);


    /* @ngInject */
    function alertModal($mdDialog) {
        
        return {
            open: open
        }
        
        function open(title, ev, content) {
            var alert = $mdDialog.alert()
                    .clickOutsideToClose(true)
                    .title(title)
                    .ok('确定')
                    .targetEvent(ev)
            if (content) {
                alert.htmlContent(content);
            }
            var dialog = $mdDialog.show(alert);
            return dialog;
        }
       
    }

})();
