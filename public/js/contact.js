function submitData(){
    let name = document.getElementById('name').value
    let email = document.getElementById('email').value
    let phone = document.getElementById('phone').value
    let subject = document.getElementById('subject').value
    let message = document.getElementById('message').value

    // console.log(name)
    // console.log(email)
    // console.log(phone)
    // console.log(subject)
    // console.log(message)

    if(name == ''){
        return alert('Name wajib di isi');
    } else if(email == ''){
        return alert('Email wajib di isi');
    } else if (phone == '') {
        return alert('Phone wajib di isi');
    } else if (subject == '') {
        return alert('Subject wajib di isi');
    } else if (message == '') {
        return alert('Message wajib di isi');
    }

    let receiver = 'maudy@gmail.com'

    let mail = document.createElement('a')
    mail.href = `mailto:${receiver}?subject=${subject}&body=Hello my name ${name}, ${message}, let's talk with me asap ${phone}`
    mail.click()
}