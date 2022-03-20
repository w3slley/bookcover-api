function find(str, term, startsBy=0){
    let len = 0
    let pos = null
    for(let i=startsBy; i < str.length; i++){
        if(str[i] == term[len])
            len++
        else
            len=0
        if(len == term.length){
            pos = i + 1 - term.length //gets position i-term.length but has to add 1 given that startsBy has default value 0
            break
        }
    }
    if(pos != null) return pos

    return -1
}

module.exports.getLinkGoogle = function (data){
    let init = find(data, 'https://www.goodreads.com/book/show/')
    let final = find(data, "&", init+10)
    let linkGoogle = data.slice(init, final)
    return linkGoogle
}
module.exports.getLinkGoodreads = function (data){
    console.log(data)
    let init = find(data, '<img src="https://i.gr-assets.com/images/')
    let final = find(data, '"', init+10)
    let linkGoodreads = data.slice(init+10, final)
    return linkGoodreads
}
