function sendPostRequest(req_data, req_url) {
  $.ajax({
    url: req_url,
    dataType: 'json',
    type: 'POST',
    contentType: 'application/json',
    data: JSON.stringify(req_data),
    processData: false,
    success: function(data, textStatus, jQxhr){
      console.log('OK');
      console.log(data);
      $("#successModalMessage").text(data.message);
      $("#successModal").modal();
    },
    error: function(jqXhr, textStatus, errorThrown){
      var _clientmsg = 'Client: ' + errorThrown;
      var _serverJSON = $.parseJSON(jqXhr.responseText);
      var _servermsg = 'Server: ' + _serverJSON.message;
      $("#errorModalMessageClient").text(_clientmsg);
      console.log(_clientmsg);
      $("#errorModalMessageServer").text(_servermsg);
      $("#errorModal").modal();
    }
  });
}

function confirmRemoveNode(_uuid) {
  var modal_message = 'Are you sure you want to remove this node?';
  $("#confirmModalMessage").text(modal_message);
  $('#confirm_action').click(function() {
    $('#confirmModal').modal('hide');
    removeNode(_uuid);
  });
  $("#confirmModal").modal();
}

function confirmRemoveNodes(_uuids) {
  var modal_message = 'Are you sure you want to remove ' + _uuids.length + ' node(s)?';
  $("#confirmModalMessage").text(modal_message);
  $('#confirm_action').click(function() {
    $('#confirmModal').modal('hide');
    removeNodes(_uuids);
  });
  $("#confirmModal").modal();
}

function removeNode(_uuid) {
  var _csrftoken = $("#csrftoken").val();
  
  var _url = '/action/' + _uuid;
  var data = {
      csrftoken: _csrftoken,
      action: 'delete'
  };
  sendPostRequest(data, _url);
}

function removeNodes(_uuids) {
  var _csrftoken = $("#csrftoken").val();
  
  var _url = '/actions';
  var data = {
      csrftoken: _csrftoken,
      uuids: _uuids, 
      action: 'delete'
  };
  sendPostRequest(data, _url);
}

function nodesView(context) {
  window.location.href = '/context/' + context + '/active';
}

function refreshCurrentNode() {
  location.reload();
}

function refreshTableNow(table_id) {
  var table = $('#' + table_id).DataTable();
  table.ajax.reload();
}

function queryNode(_uuid) {
  var _csrftoken = $("#csrftoken").val();
  var _query = $("#query").val();
  
  var _url = '/query/run';
  var data = {
      csrftoken: _csrftoken,
      context: "",
      platform: "",
      uuid_list: [_uuid],
      host_list: [],
      query: _query,
      repeat: 0
  };
  sendPostRequest(data, _url);
}

function queryNodes(_uuids) {
  var _csrftoken = $("#csrftoken").val();
  var _query = $("#query").val();
  
  var _url = '/query/run';
  var data = {
      csrftoken: _csrftoken,
      context: "",
      platform: "",
      uuid_list: _uuids,
      host_list: [],
      query: _query,
      repeat: 0
  };
  sendPostRequest(data, _url);
}

function showQueryNode(_uuid) {
  $('#query_action').click(function() {
    $('#queryModal').modal('hide');
    queryNode(_uuid);
  });
  $("#queryModal").modal();
}

function showQueryNodes(_uuids) {
  $('#query_action').click(function() {
    $('#queryModal').modal('hide');
    queryNodes(_uuids);
  });
  $("#queryModal").modal();
}