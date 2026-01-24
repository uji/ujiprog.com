document.getElementById('copy-link-btn').addEventListener('click', function() {
  var btn = this;
  var url = btn.getAttribute('data-url');
  var linkIcon = btn.querySelector('.link-icon');
  var checkIcon = btn.querySelector('.check-icon');

  navigator.clipboard.writeText(url).then(function() {
    btn.classList.add('copied');
    linkIcon.classList.add('hidden');
    checkIcon.classList.remove('hidden');

    setTimeout(function() {
      btn.classList.remove('copied');
      linkIcon.classList.remove('hidden');
      checkIcon.classList.add('hidden');
    }, 2000);
  });
});
