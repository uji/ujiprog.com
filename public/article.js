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

// Add heading anchor links
document.addEventListener('DOMContentLoaded', function() {
  var articleContent = document.querySelector('.article-content');
  if (!articleContent) return;

  var headings = articleContent.querySelectorAll('h2[id]');

  headings.forEach(function(heading) {
    var id = heading.getAttribute('id');
    if (!id) return;

    // Create anchor icon (positioned to the left)
    var anchor = document.createElement('a');
    anchor.href = '#' + id;
    anchor.className = 'heading-anchor';
    anchor.setAttribute('aria-label', 'Link to this heading');
    anchor.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>';

    // Add click handler to anchor icon
    anchor.addEventListener('click', function(e) {
      e.preventDefault();
      var url = window.location.origin + window.location.pathname + '#' + id;
      navigator.clipboard.writeText(url).then(function() {
        history.pushState(null, null, '#' + id);
      });
      heading.scrollIntoView({ behavior: 'smooth', block: 'start' });
    });

    // Create link that wraps the heading text
    var textLink = document.createElement('a');
    textLink.href = '#' + id;
    textLink.className = 'heading-link';

    // Move heading content into the text link
    while (heading.firstChild) {
      textLink.appendChild(heading.firstChild);
    }

    // Add click handler to text link
    textLink.addEventListener('click', function(e) {
      e.preventDefault();
      var url = window.location.origin + window.location.pathname + '#' + id;
      navigator.clipboard.writeText(url).then(function() {
        history.pushState(null, null, '#' + id);
      });
      heading.scrollIntoView({ behavior: 'smooth', block: 'start' });
    });

    // Insert anchor first, then text link
    heading.appendChild(anchor);
    heading.appendChild(textLink);
  });
});
