(function () {
    'use strict';
    angular.module('app')
        .controller('LogLinksCtrl', LogLinksCtrl);
    /* @ngInject */
    function LogLinksCtrl(logBackend, $stateParams) {
        var tempLogQuery = {
            cluster: $stateParams.cluster || '',
            user: $stateParams.user || '',
            from: $stateParams.from,
            to: $stateParams.to,
            app: $stateParams.app,
            task: $stateParams.task,
            path: $stateParams.path,
            page: 1,
            size: 50
        };
        var self = this;

        //Infinite Scroll setting
        self.infiniteLogs = {
            numLoaded_: 0,
            toLoad_: 0,
            items: [],

            // Required.
            getItemAtIndex: function (index) {
                if (index > this.numLoaded_) {
                    this.fetchMoreItems_(index);
                    return null;
                }
                return this.items[index];
            },

            // Required.
            getLength: function () {
                return this.numLoaded_ + 5;
            },

            fetchMoreItems_: function (index) {
                if (this.toLoad_ < index) {
                    this.toLoad_ += 50;
                    logBackend.searchLogs(tempLogQuery).get(angular.bind(this, function (data) {
                        if (!data.data.results) {
                            this.numLoaded_ = this.toLoad_ - 50;
                        } else {
                            tempLogQuery.page += 1;
                            this.items = this.items.concat(data.data.results);
                            this.numLoaded_ = this.toLoad_;
                        }
                    }));
                }
            }
        };

        activate();

        function activate() {

        }
    }
})();
