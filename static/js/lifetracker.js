$(document).ready(function(){
   var $form = $('#emailForm');
   $form.submit(function(){
      $.ajax({
        type: "POST",
        url: $(this).attr('action'),
        data: $(this).serialize(),
        success:  function(response){
          console.log ("LOLOL");
          $("#emailFormFinished").html("<strong> Thanks! Email saved! </strong>");
        }
      });
      return false;
   });
});