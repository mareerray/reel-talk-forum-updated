document.addEventListener('DOMContentLoaded', function () {
    const tabButtons = document.querySelectorAll('.nav-link');
    const tabContents = document.querySelectorAll('.tab-pane');

    tabButtons.forEach(btn => {
        btn.addEventListener('click', function () {
        // Remove active from all
        tabButtons.forEach(b => b.classList.remove('active'));
        tabContents.forEach(c => c.classList.remove('active'));

        // Add active to clicked button and show correct tab
        this.classList.add('active');
        const tabId = this.getAttribute('data-tab');
        document.getElementById(tabId).classList.add('active');
        });

        // If this is the profile tab, load profile info
        if (this.id === 'my-profile-button') {
            if (typeof loadUserProfile === 'function') {
            loadUserProfile();
            }
        }
    });
});


// document.addEventListener("DOMContentLoaded", function() {
//     // console.log("tabview.js loaded");
    
//     // Target the correct tab container ID (tabView instead of profileTabs)
//     const tabContainer = document.getElementById("tabView");
    
//     if (!tabContainer) {
//         console.error("Tab container #tabView not found!");
//         return;
//     }
    
//     // console.log("Tab container found:", tabContainer);
    
//     // Restore last active tab
//     const activeTab = localStorage.getItem("activeProfileTab");
//     // console.log("Stored active tab:", activeTab);
    
//     if (activeTab) {
//         let tabElement = document.querySelector(`[data-tab="${activeTab}"]`);
//         // console.log("Found tab element for stored tab:", tabElement);
        
//         if (tabElement) {
//             // console.log("Clicking stored tab");
//             tabElement.click(); // Switch to the saved tab
//         }
//     }

//     // Save active tab on click - use the correct selector for your HTML
//     document.querySelectorAll("#tabView button").forEach(tab => {
//         // console.log("Adding click listener to tab:", tab.id);
        
//         tab.addEventListener("click", function() {
//             const target = this.getAttribute("data-tab");
//             console.log("Tab clicked, saving target:", target);
//             localStorage.setItem("activeProfileTab", target);
//         });
//     });
// });
