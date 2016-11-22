/**
 * Created by my9074 on 2016/11/22.
 */
(function () {
    'use strict';
    angular.module('app')
        .filter('num', num);

    /* @ngInject */
    function num() {
        //////
        return function(input) {
            return parseFloat(input);
        };
    }
})();
