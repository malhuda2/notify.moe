window.gravatarAvailable = function(available) {
	var gravatarText = document.getElementById('gravatar-text');
	gravatarText.innerHTML = 'Add a gravatar image';
	gravatarText.className = available ? 'finished' : 'not-finished';
};