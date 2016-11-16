(function () {
    'use strict';
    angular.module('app')
        .config(configure);

    /* @ngInject */
    function configure($locationProvider, $urlRouterProvider, $interpolateProvider, NotificationProvider, $stateProvider) {

        $stateProvider
            .state('home', {
                url: '',
                templateUrl: '/src/app.html',
                controller: 'RootCtrl as vm',
                abstract: true
            })
            .state('home.appmonitor', {
                url: '/appmonitor',
                templateUrl: '/src/monitor/app/list/list.html',
                controller: 'ListAppCtrl as vm'
            })
            .state('home.instancemonitor', {
                url: '/instancemonitor',
                templateUrl: '/src/monitor/instance/list/list.html'
            });

        $locationProvider.html5Mode(true);
        //warning: otherwise(url) will be redirect loop on state with errored resolve
        $urlRouterProvider.otherwise(function ($injector) {
            var $state = $injector.get('$state');
            $state.go('home.appmonitor');
        });
        $interpolateProvider.startSymbol('{/');
        $interpolateProvider.endSymbol('/}');

        NotificationProvider.setOptions({
            delay: 3000 * 3,
            positionX: 'right',
            positionY: 'top',
            replaceMessage: true,
            startTop: 20,
            startRight: 260
        });
    }
})();
