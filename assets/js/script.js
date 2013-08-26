function Controller() {
}

Controller.prototype.get_structure = function(callback) {
    var self = this;
    $.getJSON('/content.json', function(data) {
        self.structure = data;
        callback();
    });
};

Controller.prototype.populate_toc = function() {
    var self = this;
    var populate_obj = function(obj, div) {
        if (obj.children && obj.children[0] && obj.children[0].children) {
            var child = $('<div></div>');
            var id = obj.id ? obj.id : obj.name;
            child.append('<a class="accordion-toggle" data-toggle="collapse" href="#' + id + '">' + obj.name + '</a>');
            var children = $('<div class="indent toc_list collapse" id="' + id + '"></div>');
            $.each(obj.children, function(index, item) {
                populate_obj(item, children);
            });
            child.append(children);
            div.append(child);
        } else {
            var clickable = $('<a href="#">' + obj.name + '</a><br/>');
            self.bind_click(clickable, obj.children);
            div.append(clickable);
        }
    }
    $.each(self.structure, function(index, item) {
        populate_obj(item, $('#toc'));
    });
};

Controller.prototype.bind_click = function(clickable, content) {
    var cnt = $('#content');
    clickable.click(function() {
        cnt.html('');
        $.each(content, function(index, obj) {
            var row = $('<div class="row"></div>');
            row.append('<img class="col-md-12" src="/content/' + obj.filename + '"/>');
            row.append('<h5>' + obj.name + '</h5>');
            row.append("<hr/>");
            cnt.append(row);
        });
    });
};

$(window).load(function() {
    ctl = new Controller();
    ctl.get_structure(function() {
        ctl.populate_toc();
        $('#toc > a:first').trigger('click');
    });
});
