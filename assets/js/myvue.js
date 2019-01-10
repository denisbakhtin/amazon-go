//obsolete 
$(document).ready(function () {

  Vue.config({
    delimiters: ["(", ")"]
  });

  var variations_exist = $('#variations').length > 0;
  var navbar_cart_exist = $('#navbar-cart').length > 0;
  var navbar_account_exist = $('#navbar-account').length > 0;
  var registration_exist = $('#modal-login').length > 0;
  var tag_edit_form_exist = $('#tag-edit-form').length > 0;
  var visible_rows = 7;
  //superagent request object
  var request = window.superagent;

  //write to dom new variation data
  function apply_variation(body) {
    if (body.Title.Valid) {
      $("h1.product-name span").text(body.Title.String);
    }
    if (body.Feature.Valid) {
      var feature_text = "<ul><li>" + body.Feature.String.split("<br/>").join("</li><li>") + "</li></ul>";
      $(".short-description div").html(feature_text);
    }
    if (body.ImageArray.length > 0) {
      //var image_array = body.Images.String.replace("\}","").replace("\{","").split(",");
      $("#current-image-box > img").attr("src", body.ImageArray[0]);
      $("#current-image-box .more-views ul li").remove();
      for (i = 0; i < body.ImageArray.length; i++) {
        $('<li class="image-panel panel panel-default"><div class="product-image-preview"><img src="' + body.ImageArray[i] + '" title="Click to View"></div></li>').appendTo("#current-image-box .more-views ul");
      }
      $("#current-image-box .more-views ul li:first-child").addClass("active");
    }
  }

  //+++++++++++++++++++++++++++ #variations VM +++++++++++++++++++++++++
  // check existance of dom element, otherwise vue js fails to initialize
  if (variations_exist) {
    // register the grid component
    // use `replace: true` because without the wrapping <table>
    // the template won't be parsed properly
    Vue.component('grid', {
      template: '#grid-template',
      replace: true,
      created: function () {
        this.ascending = {}
      },
      methods: {
        sortBy: function (key) {
          var asc = this.ascending[key] = !this.ascending[key]
          this.data.sort(function (a, b) {
            var res = a[key] > b[key]
            if (asc) res = !res
            return res ? 1 : -1
          })
        },
        setvariation: function (index) {
          for (i = 0; i < data.length; i++) {
            data[i]['html_class'] = '';
          }
          data[index]['html_class'] = 'success';
          var url = "/variations/" + parseInt(data[index]['id']);
          request.get(url, function (response) {
            if (response.ok) {
              apply_variation(response.body);
            }
          });
        },
        add_to_cart: function (index) {
          var url = "/cart/add";
          request
            .post(url)
            .type('form') //mimic html form submit, otherwise gorilla scheme won't decode parameters
            .send({
              asin: data[index]["asin"],
              quantity: "1"
            })
            .end(onResponseAdd);

          function onResponseAdd(err, res) {
            console.log("Res status: ", res.status);
            if (res.ok) {
              $("#navbar-cart .count").text(res.body.Quantity);
            }
          }
        },
        filter: function (key) {
          var value = $("#select-" + key).val();
          var defaultvalue = $("#select-" + key + " option:first-child").val();
          if (value != defaultvalue) {
            for (i = 0; i < data.length; i++) {
              if (data[i][key] != value) {
                data[i]['visible'] = false;
              } else {
                data[i]['visible'] = true;
              }
            }
            //hide buttons when filter is applied
            variations.$data.gridOptions.show_more = false;
            variations.$data.gridOptions.show_less = false;
            //reset other select filters
            //this is not vue-idiomatic, but... still
            for (i = 0; i < columns.length; i++) {
              if (columns[i]['key'] != key) {
                $("#select-" + columns[i]['key']).val(columns[i]['title']);
              }
            }
          } else {
            for (i = 0; i < data.length; i++) {
              if (i < visible_rows) {
                data[i]['visible'] = true;
              } else {
                data[i]['visible'] = false;
              }
            }
            //hide buttons when filter is applied
            variations.$data.gridOptions.show_more = data.length > visible_rows ? true : false;
            variations.$data.gridOptions.show_less = false;
          }

        },
        showall: function () {
          for (i = visible_rows; i < data.length; i++) {
            data[i]['visible'] = true;
            $("#variation-" + i).fadeIn("slow");
          }
          variations.$data.gridOptions.show_more = false;
          variations.$data.gridOptions.show_less = true;
        },
        showless: function () {
          for (i = visible_rows; i < data.length; i++) {
            data[i]['visible'] = false;
            $("#variation-" + i).fadeOut("slow");
          }
          variations.$data.gridOptions.show_more = true;
          variations.$data.gridOptions.show_less = false;
        }
      }
    });

    function initialize_columns() {
      var result = []
      $("#variations thead th").each(function (index) {
        result[index] = {
          key: $(this).attr('data-key'),
          title: $(this).text(),
          uniquevalues: []
        };
        var values = {};
        for (i = 0; i < data.length; i++) {
          values[data[i][result[index]['key']]] = 1;
        }
        var tmparray = Object.keys(values);
        tmparray.unshift(result[index]['title']);
        result[index]['uniquevalues'] = tmparray;
      });
      return result;
    };

    function initialize_data() {
      var result = []
      $("#variations tbody tr").each(function (index_tr) {
        var row = {};
        $(this).children('td').each(function (index_td) {
          row[$(this).attr('data-key')] = $(this).text();
        });
        //row visibility
        if (index_tr < visible_rows) {
          row['visible'] = true;
        } else {
          row['visible'] = false;
        }
        result[index_tr] = row;
      });

      return result;
    };

    var data = initialize_data();
    var columns = initialize_columns();
    if (data.length > 0) {
      data[0]['html_class'] = 'success'
    }
    var variations = new Vue({
      el: '#variations',
      data: {
        gridOptions: {
          data: data,
          columns: columns,
          show_more: data.length > visible_rows ? true : false,
          show_less: false
        }
      }
    });
  }

  //+++++++++++++++++++++++++++ #navbar-cart VM +++++++++++++++++++++++++
  // check existance of dom element, otherwise vue js fails to initialize
  if (navbar_cart_exist) {
    var cart = new Vue({
      el: "#navbar-cart",
      data: {
        quantity: '',
      }
    });

    request
      .get('/cart/get')
      .end(onResponseCart);

    function onResponseCart(err, res) {
      console.log("Res status: ", res.status);
      if (res.ok) {
        cart.$data.quantity = res.body.Quantity;
      }
    }
  }

  //+++++++++++++++++++++++++++ #navbar account VM +++++++++++++++++++++++++
  // check existance of dom element, otherwise vue js fails to initialize
  if (navbar_account_exist) {
    var account = new Vue({
      el: "#navbar-account",
      data: {
        email: 'Account',
        signed_in: false
      },
      methods: {
        logoff: function () {
          request
            .post('/sign_out')
            .end(onResponseLogoff);

          function onResponseLogoff(err, res) {
            console.log("Resss status: ", res.status);
            if (res.ok) {
              //account.$data.email = "Account";
              //account.$data.signed_in = false;
              location.reload();
            }
          }
        }
      }
    });

    request
      .get('/signed_in')
      .end(onResponseSigned);

    function onResponseSigned(err, res) {
      console.log("Res status: ", res.status);
      if (res.ok && res.body.Email != undefined && res.body.Email.length > 0) {
        account.$data.email = res.body.Email;
        account.$data.signed_in = true;
      }
    }
  }

  //+++++++++++++++++++++++++++ #registration/sign-in VM +++++++++++++++++++++++++
  // check existance of dom element, otherwise vue js fails to initialize
  if (registration_exist) {
    var reg = new Vue({
      el: "#modal-login",
      data: {
        email: '',
        password: '',
        password_confirmation: '',
        error_text: ''
      },
      methods: {
        signin: function () {
          request
            .post('/sign_in')
            .type('form')
            .send({
              email: reg.$data.email,
              password: reg.$data.password
            })
            .end(onResponseSignin);

          function onResponseSignin(err, res) {
            console.log("Res status: ", res);
            if (res.ok) {
              //account.$data.email = reg.$data.email;
              //account.$data.signed_in = true;
              $("#modal-login").modal('hide');
              location.reload();
            } else {
              reg.$data.error_text = res.text;
            }
          }
        },
        clear_error: function () {
          reg.$data.error_text = '';
        },
        signup: function () {
          request
            .post('/sign_up')
            .type('form')
            .send({
              email: reg.$data.email,
              password: reg.$data.password,
              password_confirmation: reg.$data.password_confirmation
            })
            .end(onResponseSignup);

          function onResponseSignup(err, res) {
            console.log("Res status: ", res);
            if (res.ok) {
              //account.$data.email = reg.$data.email;
              //account.$data.signed_in = true;
              $("#modal-login").modal('hide');
              location.reload();
            } else {
              reg.$data.error_text = res.text;
            }
          }
        }
      }
    });
  }

  //+++++++++++++++++++++++++++ #tag edit form VM +++++++++++++++++++++++++
  // check existance of dom element, otherwise vue js fails to initialize
  if (tag_edit_form_exist) {
    var tag_form = new Vue({
      el: "#tag-edit-form",
      data: {
        tags_and_products: [],
        changes_visible: false,
        title_includes: $("#tag-edit-form inputTitleIncludes").val(),
        title_excludes: $("#tag-edit-form inputTitleExcludes").val(),
        total_products: 0,
      },
      methods: {
        viewChanges: function () {
          request
            .post(tag_form.$el.action + '/view_mask.json')
            .type('form')
            .send({
              title_includes: tag_form.$data.title_includes,
              title_excludes: tag_form.$data.title_excludes
            })
            .end(onResponseView);

          function onResponseView(err, res) {
            console.log("Res status: ", res);
            if (res.ok) {
              tag_form.$data.tags_and_products = res.body;
              tag_form.$data.changes_visible = true;
              var total_products = 0;
              for (i = 0; i < tag_form.$data.tags_and_products.length; i++) {
                total_products = total_products + tag_form.$data.tags_and_products[i].ProductCount;
              }
              tag_form.$data.total_products = total_products;
            } else {
              //reg.$data.error_text = res.text;
            }
          }
        },
        applyChanges: function () {
          request
            .post(tag_form.$el.action + '/apply_mask.json')
            .type('form')
            .send({
              title_includes: tag_form.$data.title_includes,
              title_excludes: tag_form.$data.title_excludes
            })
            .end(onResponseApply);

          function onResponseApply(err, res) {
            console.log("Res status: ", res);
            if (res.ok) {
              tag_form.$data.tags_and_products = [];
              tag_form.$data.changes_visible = false;
              tag_form.$data.total_products = 0;
            } else {
              alert(res.text);
            }
          }
        }
      }
    });
  }

});