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
            var id = obj.id ? obj.id : obj.name;
            var link = $('<a class="accordion-toggle" data-toggle="collapse" href="#' + id + '"><span>' + obj.name + '</span></a>');
            var children = $('<div class="toc_list collapse" id="' + id + '"></div>');
            children.on('shown.bs.collapse', function() {
                link.addClass('dir-in');
            });
            children.on('hidden.bs.collapse', function() {
                link.removeClass('dir-in');
            });
            $.each(obj.children, function(index, item) {
                populate_obj(item, children);
            });
            var container = $('<div></div>');
            container.append(link);
            container.append(children);
            div.append(container);
        } else {
            var clickable = $('<a href="#"><span>' + obj.name + '</span></a><br/>');
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
        $('.link-in').removeClass('link-in');
        clickable.addClass('link-in');
        cnt.html('');
        $.each(content, function(index, obj) {
            var row = $('<div class="content_item"></div>');
            row.append('<img class="content_item_img" src="/content/' + obj.filename + '" alt="' + obj.name + '"/>');
            cnt.append(row);
        });
    });
};

function init() {
    $(window).resize(function() {
        $('#toc').height($('.sidebar').height() - $('.sidebar').width() - $('#toc').position().top - parseInt($('.sidebar').css("margin-top")));
    });
    $(window).load(function() {
        ctl = new Controller();
        ctl.get_structure(function() {
            ctl.populate_toc();
            $(window).trigger('resize');
            $('#toc > a:first').trigger('click');
        });
    });
}

init();
