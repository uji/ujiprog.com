var icons = {
  zenn: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M.264 23.771h4.984c.264 0 .498-.147.645-.352L19.614.874c.176-.293-.029-.645-.381-.645h-4.72c-.235 0-.44.117-.557.323L.03 23.126c-.088.176.029.645.234.645zM17.445 23.419l6.479-10.408c.205-.323-.029-.733-.41-.733h-4.691c-.176 0-.352.088-.44.235l-6.655 10.643c-.176.264.029.616.352.616h4.926c.176 0 .352-.088.44-.353z" fill="#3EA8FF"/></svg>',
  note: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><rect x="2" y="2" width="20" height="20" rx="5" fill="#FFFFFF"/><text x="12" y="17" text-anchor="middle" fill="#000000" font-family="Arial" font-size="14" font-weight="bold">n</text></svg>',
  speakerdeck: '<svg width="16" height="16" viewBox="41 25 32 20" xmlns="http://www.w3.org/2000/svg"><path d="M54.3665414,37.5 L47.25,37.5 C43.7982203,37.5 41,34.7017797 41,31.25 C41,27.7982203 43.7982203,25 47.25,25 L55.5526316,25 C56.9333435,25 58.0526316,26.1192881 58.0526316,27.5 C58.0526316,28.8807119 56.9333435,30 55.5526316,30 L47.1221805,30 C46.4318245,30 45.8721805,30.5596441 45.8721805,31.25 C45.8721805,31.9403559 46.4318245,32.5 47.1221805,32.5 L54.2387218,32.5 C57.6905015,32.5 60.4887218,35.2982203 60.4887218,38.75 C60.4887218,42.2017797 57.6905015,45 54.2387218,45 L43.5,45 C42.1192881,45 41,43.8807119 41,42.5 C41,41.1192881 42.1192881,40 43.5,40 L54.3665414,40 C55.0568973,40 55.6165414,39.4403559 55.6165414,38.75 C55.6165414,38.0596441 55.0568973,37.5 54.3665414,37.5 Z M59.6267041,45 C61.2891288,43.8757084 62.4773068,42.0834962 62.8209549,40 L66.8554291,40 C67.5341396,40 68.0843433,39.4403559 68.0843433,38.75 L68.0843433,31.25 C68.0843433,30.5596441 67.5341396,30 66.8554291,30 L59.5263158,30 C60.1100991,29.3365544 60.4650753,28.460443 60.4650753,27.5 C60.4650753,26.539557 60.1100991,25.6634456 59.5263158,25 L68.0843433,25 C70.7991855,25 73,27.2385763 73,30 L73,40 C73,42.7614237 70.7991855,45 68.0843433,45 L59.6267041,45 Z" fill="#009287"/></svg>',
  blog: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M19 3H5C3.9 3 3 3.9 3 5V19C3 20.1 3.9 21 5 21H19C20.1 21 21 20.1 21 19V5C21 3.9 20.1 3 19 3ZM7 7H17V9H7V7ZM7 11H17V13H7V11ZM7 15H14V17H7V15Z" fill="#4A4B4A"/></svg>'
};

fetch('/articles.json')
  .then(function(res) { return res.json(); })
  .then(function(data) {
    var container = document.getElementById('zenn-articles');
    data.articles.forEach(function(article) {
      var a = document.createElement('a');
      a.href = article.url;
      a.className = 'article-card';
      var date = new Date(article.published_at).toLocaleDateString('ja-JP');
      var icon = icons[article.platform] || icons.zenn;
      a.innerHTML = icon + '<span class="article-title">' + article.title.replace(/</g, '&lt;').replace(/>/g, '&gt;') + '</span><span class="article-date">' + date + '</span>';
      if (article.platform !== 'blog') {
        a.target = '_blank';
        a.rel = 'noopener noreferrer';
      }
      container.appendChild(a);
    });
  });
