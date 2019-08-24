//SPDX-License-Identifier: Apache-2.0

var car = require('./controller.js');

module.exports = function(app){

  app.get('/get_car/:id', function(req, res){
    car.get_car(req, res);
  });
  app.get('/add_car/:car', function(req, res){
    car.add_car(req, res);
  });
  app.get('/get_all_cars', function(req, res){
    car.get_all_cars(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    car.change_holder(req, res);
  });
}
