const wrapper = document.querySelector('.wrapper');
const loginLink = document.querySelector('.login-link');
const registerLink = document.querySelector('.register-link');
const buttonLoginPopup = document.querySelector('.button__login-popup');
const buttonRegisterPopup = document.querySelector('.button__register-popup');
const closeButton = document.querySelector('.ri-close-line');

registerLink.addEventListener('click', () => {
    wrapper.classList.add('active');
});

loginLink.addEventListener('click', () => {
    wrapper.classList.remove('active');
});

buttonLoginPopup.addEventListener('click', () => {
    wrapper.classList.add('login-popup');
});

closeButton.addEventListener('click', () => {
    wrapper.classList.remove('login-popup');
    wrapper.classList.remove('active');
});
                    //////// notif ////////
        const notif = document.querySelector('.notifConnected');

        const notifBubble = document.querySelector('.notifBubble');
        // const notifBubbleAnim = gsap.to(notifBubble, {transform: "scale(1.2)", duration: 0.25, ease: "none.out", yoyo: true, repeat: -1, paused: true})

        const notifBubbleWave = document.querySelector('.notifBubbleWave');
        // const notifBubbleWaveAnim = gsap.to(notifBubbleWave, {transform: "scale(1.7)", duration: 0.5, opacity: 0, ease: "none.out", repeat: -1, paused: true})

                    //////// NavBar ////////
        const logo = document.querySelectorAll('.logo');
        const logoAnim = gsap.to(logo, {transform: "scale(1.5)", filter: "blur(0px)", duration: 1.9, ease: "none.out", paused: true});

        const navLinks = document.querySelectorAll('.nav__links');
        const navLinksAnim = gsap.from(navLinks, {top: 30, opacity: 0, duration: 1.9, stagger: 0.5, ease: "none.out", delay: 1, paused: true});

        const filterDropdown = document.querySelectorAll('.filter');
        const filterDropdownAnim = gsap.fromTo(filterDropdown, {scale: 0}, {scale: 1, duration: 1.5, paused: true, ease: "bounce.out", delay: 3});

        const TL = gsap.timeline({paused: true});
            TL
            .to(notif, {left: 1400, opacity: 1, delay: 2, duration: 2, ease: "bounce.out"})
            .to(notif, {left: 1500, opacity: 1, duration: 2, yoyo: true, repeat: 2, ease: "bounce.in"})
            .to(notif, {left: 1870, opacity: 1, duration: 1, ease: "none.out"})

        const TL2 = gsap.timeline({paused: true});
            TL2
            .to(notifBubble, {scale: 70, duration: 1, ease: "none.out"})
            .to(notifBubble, {scale: 90, duration: 0.25, ease: "none.out", yoyo: true, repeat: 38, delay: 1})
            .to(notifBubble, {scale: 0, duration: 1, ease: "none.out"});

        const TL3 = gsap.timeline({paused: true});
            TL3
            .to(notifBubbleWave, {scale: 70, duration: 1, ease: "none.out"})
            .to(notifBubbleWave, {scale: 200, duration: 0.5, delay: 1, opacity: 0, ease: "none.out", repeat: 19})
            .to(notifBubbleWave, {scale: 0, duration: 1, ease: "none.in"});

        window.addEventListener('load', () => {
            notifs();
            navLinksAnim.play();
            logoAnim.play();
            filterDropdownAnim.play();
        })

        function notifs() {
            TL2.play();
            TL3.play();
            TL.play();
        }

        navLinks.forEach(link => {
            link.addEventListener('click', function() {
                gsap.to(link, {top: 20, duration: 1, ease: "power4.out"})
            });

            link.addEventListener('mouseout', () => {
                gsap.to(link, {top: 0, duration: 1, ease: "power4.out"})
            });
        });

                    ////// like ////////
        // const img2All = document.querySelectorAll('.img2');
        const img2 = document.querySelector('.img2');
        const img4 = document.querySelector('.img4');
        const like = gsap.to(img2, {transform: "scale(100%)", opacity: 1, duration: 1, ease: "power4.out", paused: true})
        const unlike = gsap.to(img4, {transform: "scale(100%)", opacity: 1, duration: 1, ease: "power4.out", paused: true})

        let liked = false;
        let unliked = false;

        img2.addEventListener('click', function() {
            if (!liked) {
                like.play();
                liked = true;
            } else {
                like.reverse();
                liked = false;
            }
        })

        img4.addEventListener('click', function() {
            if (!unliked) {
                unlike.play();
                unliked = true;
            } else {
                unlike.reverse();
                unliked = false;
            }
        })