'use strict';

/* Controllers */

angular.module('myApp.controllers').
  controller('GameSelectionCtrl', function ($scope, response) {
    $scope.data = response.data;
  });
