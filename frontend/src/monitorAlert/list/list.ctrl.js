/**
 * Created by lixiaoyan on 2017/3/24.
 */
(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorAlertCtrl', MonitorAlertCtrl);
    /* @ngInject */

    function MonitorAlertCtrl(rules,monitorAlertBackend, confirmModal, Notification, $state) {
        var self = this;

        self.rules = rules.data.rules;
        // self.ALERT_STATUS = {
        //     Enabled: '激活',
        //     Disabled: '未激活'
        // };

        self.deleteAlert = deleteAlert;
        self.activateAlert = activateAlert;
        activate();

        function activate() {

        }

        function deleteAlert(id, monitorClass) {
            var data = {'class' : monitorClass};

            confirmModal.open("是否确认删除该监控告警规则?").then(function () {
                monitorAlertBackend.rule(id)
                    .remove(data, function (data) {
                        Notification.success('删除成功');
                        $state.go('home.monitorAlert', null, {reload: true});
                    })
            });
        }


        function activateAlert(data) {
            var modalText = data.status === 'Enabled' ? '是否暂停该告警?' : '是否启用该告警';

            confirmModal.open(modalText).then(function () {
                data.status = data.status === 'Enabled' ? 'Disabled' : 'Enabled';
                monitorAlertBackend.rule(data.ID).update(data, function (data) {
                    Notification.success('操作成功');
                    $state.reload();
                });
            });
        }
    }
})();