document.addEventListener('DOMContentLoaded', function () {
  const track = document.getElementById('categoryTrack');

  if (!track) return;

  fetch('/api/categories')
    .then(response => {
      if (!response.ok) throw new Error('Network response was not ok');
      return response.json();
    })
    .then(categories => {
      const chips = [];

      const allChip = `
        <button type="button" class="category-chip" data-category="All">
          🎬 All Genres
        </button>
      `;
      chips.push(allChip);

      categories.forEach(category => {
        chips.push(`
          <button type="button" class="category-chip" data-category="${category.name}">
            ${category.emoji} ${category.name}
          </button>
        `);
      });

      const html = chips.join('');
      track.innerHTML = html + html;

      track.addEventListener('click', function (event) {
        const chip = event.target.closest('.category-chip');
        if (!chip) return;

        filterPostsByCategory(chip.dataset.category);
      });
    })
    .catch(error => {
      console.error('Error fetching categories:', error);
    });
});
// document.addEventListener('DOMContentLoaded', function() {
//     // Select the "All Genres" button
//     const allGenresBtn = document.querySelector('.all-movies');
//     if (allGenresBtn) {
//         allGenresBtn.addEventListener('click', function(event) {
//             // No need for preventDefault since it's a button, but you can include it for safety
//             event.preventDefault();
//             filterPostsByCategory('All');
//         });
//     }

//     // Fetch categories from the API and add category buttons as before
//     fetch('/api/categories')
//         .then(response => {
//             if (!response.ok) {
//                 throw new Error('Network response was not ok');
//             }
//             return response.json();
//         })
//         .then(categories => {
//             // Get the category list element
//             const categoryList = document.querySelector('.category-bar ul');

//             // Remove all existing category items except the first (All Genres)
//             const allGenresItem = categoryList.querySelector('li:first-child');
//             categoryList.innerHTML = '';
//             categoryList.appendChild(allGenresItem);

//             // Add each category as a button
//             categories.forEach(category => {
//                 const li = document.createElement('li');
//                 const btn = document.createElement('button');
//                 btn.type = 'button';
//                 btn.className = 'category-btn';
//                 btn.innerHTML = `${category.emoji} ${category.name}`;
//                 btn.addEventListener('click', function() {
//                     filterPostsByCategory(category.name);
//                 });
//                 li.appendChild(btn);
//                 categoryList.appendChild(li);
//             });
//         })
//         .catch(error => {
//             console.error('Error fetching categories:', error);
//         });
// });

