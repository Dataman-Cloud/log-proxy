(function () {
    'use strict';
    angular.module('app')
        .config(configure);

    /* @ngInject */
    function configure($locationProvider, $urlRouterProvider, $interpolateProvider, $stateProvider) {

        $stateProvider
            .state('home', {
                url: '',
                templateUrl: '/src/app.html',
                controller: 'RootCtrl as vm',
                abstract: true
            })
            .state('home.appmonitor', {
                url: '/appmonitor',
                templateUrl: '/src/monitor/app/list.html',
                controller: 'MonitorListAppCtrl as vm',
                resolve: {
                    apps: listApp
                }
            })
            .state('home.appmonitor.detail', {
                url: '/:appId',
                controller: 'MonitorAppDetailCtrl as vm',
                templateUrl: '/src/monitor/app/detail.html'
            })
            .state('home.instancemonitor', {
                url: '/appmonitor/:appId/instances/:taskId',
                templateUrl: '/src/monitor/instance/detail.html',
                controller: 'MonitorInstanceCtrl as vm'
            })
            .state('home.nodemonitor', {
                url: '/nodemonitor',
                templateUrl: '/src/monitor/node/node.html',
                controller: 'MonitorNodeCtrl as vm',
                resolve: {
                    nodes: listNode
                }
            })
            .state('home.log', {
                url: '/log',
                templateUrl: '/src/log/log.html',
                controller: 'LogCtrl as vm'
            });


        /* @ngInject */
        function listApp(monitorBackend) {
            return monitorBackend.listApp().get().$promise
        }

        /* @ngInject */
        function listApp(monitorBackend) {
            return monitorBackend.listApp().get().$promise
        }

        /* @ngInject */
        function listNode(monitorBackend) {
            return monitorBackend.listNode().get().$promise
        }

        $locationProvider.html5Mode(true);
        //warning: otherwise(url) will be redirect loop on state with errored resolve
        $urlRouterProvider.otherwise(function ($injector) {
            var $state = $injector.get('$state');
            $state.go('home.appmonitor');
        });
        $interpolateProvider.startSymbol('{/');
        $interpolateProvider.endSymbol('/}');
    }
})();
