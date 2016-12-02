/**
 * Created by lixiaoyan on 2016/12/2.
 */
(function () {
    'use strict';
    angular.module('app')
        .factory('alertcurd', alertcurd);

    /* @ngInject */
    function alertcurd(alertBackend, confirmModal, Notification, $state) {
        return {
            deleteAlert: deleteAlert
        };

        function deleteAlert(id) {
            confirmModal.open("是否确认删除应用?").then(function () {
                alertBackend.alert(id)
                    .delete(function (data) {
                        Notification.success('删除成功');
                        $state.go('home.alertlist', null, {reload: true});
                    })
            });
        }
    }
})();