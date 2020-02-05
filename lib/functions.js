

function find(str, term, startsBy=0){
    let x = []
    let j = 0
    let pos = null
    for(let i=startsBy; i < str.length; i++){
        if(str[i] == term[j]){
            x.push(1)
            j++
        }
        else{
            j=0
            x=[]
        }
        if(x.length == term.length){
            pos = i - (term.length-1)
            break
        }
    }
    if(pos != null){
        return pos
    }

    return -1
}

module.exports.getLinkGoogle = function (data){
    let init = find(data, 'https://www.goodreads.com/book/show/')
    let final = find(data, "&", init+10)
    let linkGoogle = data.slice(init, final)
    return linkGoogle
}
module.exports.getLinkGoodreads = function (data){
    let init = find(data, '<img src=')
    let final = find(data, '"', init+10)
    let linkGoodreads = data.slice(init+10, final)
    return linkGoodreads
}
