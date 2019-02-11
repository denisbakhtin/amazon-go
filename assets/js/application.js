import xhook from 'xhook';
import $ from 'jquery';
window.jQuery = $;
window.$ = $;

import 'popper.js';
import 'bootstrap';
import 'malihu-custom-scrollbar-plugin';
import '../scss/application.scss';
//import Siema from 'siema';

import {
  library,
  dom
} from '@fortawesome/fontawesome-svg-core'
import {
  faCog,
  faShoppingCart,
  faChevronRight,
  faChevronLeft,
  faCheck,
  faStar,
  faHome,
  faSync,
  faTimes,
  faExternalLinkAlt
} from "@fortawesome/free-solid-svg-icons";
import {
  faFrown,
  faMehBlank,
  faStar as farStar,
  faMeh,
} from "@fortawesome/free-regular-svg-icons";

library.add(faExternalLinkAlt, faMeh, faCog, faShoppingCart, faChevronRight, faChevronLeft, faCheck, faStar, farStar, faHome, faSync, faTimes, faFrown, faMehBlank);

$.fn.visibleHeight = function () {
  //see https://stackoverflow.com/questions/24768795/get-the-visible-height-of-a-div-with-jquery
  var elBottom, elTop, scrollBot, scrollTop, visibleBottom, visibleTop;
  scrollTop = $(window).scrollTop();
  scrollBot = scrollTop + $(window).height();
  elTop = this.offset().top;
  elBottom = elTop + this.outerHeight();
  visibleTop = elTop < scrollTop ? scrollTop : elTop;
  visibleBottom = elBottom > scrollBot ? scrollBot : elBottom;
  return visibleBottom - visibleTop
}

function markotherdimoptions(id, val) {
  //premark all options of other selects with 'option-muted' class
  for (var j = 0; j < document.dims[0].length - 1; j++) {
    if (j != id) {
      $('#dim-' + j + '-select option').each(function (el) {
        $(this).addClass('option-muted');
      })
    }
  }

  for (var i = 0; i < document.dims.length; i++) {
    if (document.dims[i][id] === val) {
      for (var j = 0; j < document.dims[i].length - 1; j++) {
        if (j != id) {
          $('#dim-' + j + '-select option[value="' + document.dims[i][j] + '"]').removeClass('option-muted');
        }
      }
    }
  }
}

//set all dim selects to compatible options once a muted option has been picked
function updatedimselects(id, val) {
  for (var j = 0; j < document.dims[0].length - 1; j++) {
    if (j != id) {
      if ($('#dim-' + j + '-select option[value="' + $('#dim-' + j + '-select').val() + '"]').hasClass('option-muted')) {
        for (var i = 0; i < document.dims.length; i++) {
          if (document.dims[i][id] === val) {
            for (var k = 0; k < document.dims[i].length - 1; k++) {
              $('#dim-' + k + '-select').val(document.dims[i][k]);
            }
            return;
          }
        }
      }
    }
  }
}

//mark all incompatible dim options as muted (gray)
function markdimoptions() {
  //premark all options of other selects with 'option-muted' class
  for (var j = 0; j < document.dims[0].length - 1; j++) {
    $('#dim-' + j + '-select option').each(function (el) {
      $(this).addClass('option-muted');
    })
  }

  var vals = new Array(document.dims[0].length - 1);
  for (var i = 0; i < vals.length; i++) {
    vals[i] = $('#dim-' + i + '-select').val();
  }
  for (var i = 0; i < vals.length; i++) {
    for (var j = 0; j < document.dims.length; j++) {
      var found = true;
      for (var k = 0; k < document.dims[j].length - 1; k++) {
        if (k != i)
          found = found && (document.dims[j][k] == vals[k]);
      }
      if (found)
        $('#dim-' + i + '-select option[value="' + document.dims[j][i] + '"]').removeClass('option-muted');
    }
  }
}

function updateselectedasin() {
  var vals = new Array(document.dims[0].length - 1);
  for (var i = 0; i < vals.length; i++) {
    vals[i] = $('#dim-' + i + '-select').val();
  }
  for (var i = 0; i < document.dims.length; i++) {
    var found = true;
    for (var j = 0; j < document.dims[i].length - 1; j++) {
      found = found && (vals[j] === document.dims[i][j]);
    }
    if (found) {
      document.selected_asin = document.dims[i][document.dims[i].length - 1]
      return;
    }
  }
}

//set width, height of images in product grid
function set_image_height() {
  //set height of image box in product view
  var cw = $('.product-view #current-image-box').width();

  //set height of image box in cart view
  var cw = $('.cart-view .image a').width();
  $('.cart-view .image a').css({
    'height': (cw) + 'px'
  });
  $('.cart-view .image a').css({
    'line-height': (cw) + 'px'
  });
};

function updatevariation(asin) {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
      var obj = JSON.parse(xhttp.responseText)
      $('.offerdetails').html(obj.Offer);
      $('.product-image-box').html(obj.Images);
      $('.short-description').html(obj.Description);
    }
  };
  xhttp.open("GET", "/variations/" + asin, true);
  xhttp.send();
}

//fasten page load by appending iframe scr after onload
function set_reviews() {
  var src = $('#reviews').attr('data-src');
  var ifrm = document.createElement('iframe');
  $(ifrm).attr('src', src);
  $('#reviews').append(ifrm);
  //or add iframe html source to page body :)
  /* var id = $('#reviews').attr('data-id');
    $.get('/product_reviews/' + id, function (data) {
    console.log(data);
    $('#reviews-body').html(data);
  }); */
};

//set dynamic product and tag preview description bg color
function set_preview_bgcolor() {
  $('.tag-preview-wrapper').each(function () {
    var pbgcolor = $(this).attr('data-product-preview-bgcolor');
    var tbgcolor = $(this).attr('data-tag-preview-bgcolor');
    $(this).find('.product-preview').each(function () {
      $(this).find('.description').css("background-color", pbgcolor);
    });
    $(this).find('.tag-preview').each(function () {
      $(this).find('.description').css("background-color", tbgcolor);
    });
  })
}

$(document).ready(function () {

  //font-awesome replace <i> with svg
  dom.watch()

  //variations click
  $('#variations tbody tr').on('click', function (e) {
    $('#variations tbody tr').removeClass('bg-success');
    $(e.currentTarget).addClass('bg-success');

    updatevariation($(e.currentTarget).find('.asin').text());
  });

  //ckeditor integration
  if (document.querySelector('#ck-content')) {
    //add csrf protection to ckeditor uploads
    xhook.before(function (request) {
      if (!/^(GET|HEAD|OPTIONS|TRACE)$/i.test(request.method)) {
        request.xhr.setRequestHeader("X-CSRF-TOKEN", window.csrf_token);
      }
    });

    ClassicEditor
      .create(document.querySelector('#ck-content'), {
        language: 'en', //to set different lang include <script src="/public/js/ckeditor/build/translations/{lang}.js"></script> along with core ckeditor script
        ckfinder: {
          uploadUrl: '/admin/upload'
        },
      })
      .catch(error => {
        console.error(error);
      });
  }

  $('.dimension-select').on('change', function () {
    var id = $(this).attr('data-id');
    var val = $(this).val();

    markotherdimoptions(id, val);
    updatedimselects(id, val);
    markdimoptions();
    updateselectedasin();
    updatevariation(document.selected_asin);
  });

  //set bg colors for tag previews (mainly on home page)
  if ($('.tag-preview-wrapper').length > 0) {
    set_preview_bgcolor();
  }

  //siema carousel plugin
  /* if ($('.tag-preview-siema').length > 0) {
    $('.tag-preview-siema').each(function () {
      const id = $(this).attr('id');
      const tagPreviewSiema = new Siema({
        selector: '#' + id,
        duration: 200,
        easing: 'ease-out',
        perPage: {
          150: 1,
          450: 2,
          650: 3,
          850: 4,
          1000: 5,
          1150: 6,
        },
        startIndex: 0,
        draggable: true,
        multipleDrag: true,
        threshold: 20,
        loop: true,
        rtl: false,
        onInit: () => {},
        onChange: () => {},
      });
      $('#' + id).parent().find('.tag-product-prev').on('click', () => tagPreviewSiema.prev());
      $('#' + id).parent().find('.tag-product-next').on('click', () => tagPreviewSiema.next());
    })
  } */

  //initialize correct selected size from active variation
  set_image_height();
  set_reviews();
  if (document.dims != undefined) {
    markdimoptions();
    for (var i = 0; i < document.dims[0].length - 1; i++) {
      $('#dim-' + i + '-select').val(document.dims[0][i]);
    }
    //initialize active asin
    document.selected_asin = document.dims[0][document.dims[0].length - 1];
  }

  //get tags on category change, product edit form
  $('select#select-product-category').on('change', function (e) {
    var category_id = $(this).val();
    $('#active_category_id').val(category_id);
    $('#set_tag').submit();
  });

  //change special_price on regular_price change
  $('#regular_price input').on('change', function (e) {
    var regular_price = $(this).val() || 0;
    var discount_percent = $('#discount_percent input').val() || 0;
    $('#special_price input').val(regular_price - regular_price * discount_percent / 100);
  });
  //change special_price on discount_percent change
  $('#discount_percent input').on('change', function (e) {
    var discount_percent = $(this).val() || 0;
    var regular_price = $('#regular_price input').val() || 0;
    $('#special_price input').val(regular_price - regular_price * discount_percent / 100);
  });
  //change discount_percent on special_price change
  $('#special_price input').on('change', function (e) {
    var special_price = $(this).val() || 0;
    var regular_price = $('#regular_price input').val() || 0;
    $('#discount_percent input').val((regular_price - special_price) / regular_price * 100);
  });

  //click on category select in search form
  $('#category_search_list').on('click', 'a', function (e) {
    $('#category_search_title').text($(this).text());
    $('#category_search_id').val($(this).attr('data-search-id'));
  });

  //show category submenu on parent hover
  $('.sidebar-nav').on('mouseenter', 'ul > li', function (e) {
    var submenu = $('#' + $(this).attr('data-submenu-id'));
    var parent_offset = $(this).position();
    submenu.css({
      width: $(this).width(),
      top: -5,
      left: parent_offset.left + $(this).width() - 30
    });
    submenu.show('fast');
  });
  $('.sidebar-nav').on('mouseleave', 'ul > li', function (e) {
    var active_element = e.toElement || e.relatedTarget;
    if (!$(active_element).hasClass('submenu')) {
      $('#' + $(this).attr('data-submenu-id')).hide(0);
    }
  });
  $('.sidebar-nav').on('mouseleave', '.submenu', function (e) {
    $(this).hide(0);
  });

  $(window).resize(function () {
    set_image_height();
  });

});